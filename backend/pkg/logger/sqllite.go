package logger

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"tiny-forum/config"
	// _ "modernc.org/sqlite"
)

// DBEntry 单条日志记录（与数据库字段一一对应）
type DBEntry struct {
	Time       time.Time // RFC3339Nano UTC
	Level      string    // DEBUG / INFO / WARN / ERROR / FATAL
	Caller     string    // file:line
	Message    string    // 日志内容
	Stacktrace string    // 错误堆栈
	Fields     string    // zap 附加字段，JSON 格式
}

// ─── 数据库表结构 ────────────────────────────────────────────
//
// 按天分表：logs_20260501、logs_20260502 …
// 好处：
//   - 整表 DROP 清理过期数据，比 DELETE WHERE time<? 快得多
//   - 每张表行数可控，查询无需全表扫描
//   - 与 Makefile 的 TODAY=$(shell date +%Y%m%d) 命名完全一致
//
// 表结构（与 Makefile 中的 SELECT 字段完全匹配）：
//
//	CREATE TABLE logs_YYYYMMDD (
//	    id         INTEGER PRIMARY KEY AUTOINCREMENT,
//	    timestamp  DATETIME NOT NULL,   -- RFC3339Nano UTC，可直接字符串比较
//	    level      TEXT     NOT NULL,   -- DEBUG / INFO / WARN / ERROR / FATAL
//	    caller     TEXT,               -- file:line
//	    message    TEXT,
//	    stacktrace TEXT,
//	    fields     TEXT                -- JSON
//	);

// dbSink SQLite 异步写入器
type dbSink struct {
	db  *sql.DB
	cfg *config.DBConfig

	// 按天缓存 prepared statement，避免重复 prepare
	mu      sync.Mutex
	stmtMap map[string]*sql.Stmt // key: "YYYYMMDD"

	queue chan DBEntry
	done  chan struct{}
	wg    sync.WaitGroup
	once  sync.Once
}

var (
	globalDB   *dbSink
	globalDBMu sync.Mutex
)

// ─── 公共 API ─────────────────────────────────────────────────

// InitDB 初始化 SQLite 日志数据库。
// 可在 logger.Init 之前或之后单独调用；若在 Init 之后调用，需重新调用 Init 以接入 dbCore。
func InitDB(cfg *config.DBConfig) error {
	fmt.Printf("dsn %v", cfg.DSN)
	globalDBMu.Lock()
	defer globalDBMu.Unlock()

	if globalDB != nil {
		_ = globalDB.close()
	}

	sink, err := newDBSink(cfg)
	if err != nil {
		return err
	}
	globalDB = sink
	return nil
}

// CloseDB 等待队列写完后关闭数据库连接。
// 建议在 main 退出前调用：defer logger.CloseDB()
func CloseDB() error {
	globalDBMu.Lock()
	defer globalDBMu.Unlock()
	if globalDB == nil {
		return nil
	}
	err := globalDB.close()
	globalDB = nil
	return err
}

// QueryLogs 查询日志。
//   - level: 日志级别过滤（大写，如 "ERROR"），空字符串表示不过滤
//   - start / end: 时间范围，零值表示不限制（默认最近 24h）
//   - limit: 返回条数上限，≤0 时默认 200
//
// 自动覆盖 start~end 所跨越的所有分表，跨天时合并结果。
func QueryLogs(level string, start, end time.Time, limit int) ([]DBEntry, error) {
	globalDBMu.Lock()
	sink := globalDB
	globalDBMu.Unlock()

	if sink == nil {
		return nil, fmt.Errorf("db sink 未初始化，请先调用 InitDB")
	}
	return sink.query(level, start, end, limit)
}

// ─── 内部实现 ─────────────────────────────────────────────────

func newDBSink(cfg *config.DBConfig) (*dbSink, error) {
	if cfg.MaxBuffer <= 0 {
		cfg.MaxBuffer = 512
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 50
	}
	if cfg.FlushEvery <= 0 {
		cfg.FlushEvery = 2 * time.Second
	}

	if err := os.MkdirAll(filepath.Dir(cfg.DSN), 0755); err != nil {
		return nil, fmt.Errorf("创建 db 目录失败: %w", err)
	}

	db, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("打开 sqlite 失败: %w", err)
	}
	// 单写连接，防止 SQLITE_BUSY；读连接不受此限制
	db.SetMaxOpenConns(1)

	if _, err = db.Exec(`PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("设置 PRAGMA 失败: %w", err)
	}

	s := &dbSink{
		db:      db,
		cfg:     cfg,
		stmtMap: make(map[string]*sql.Stmt),
		queue:   make(chan DBEntry, cfg.MaxBuffer),
		done:    make(chan struct{}),
	}

	s.wg.Add(1)
	go s.loop()
	return s, nil
}

// tableKey 返回给定时间对应的分表后缀，格式 "YYYYMMDD"（UTC）
func tableKey(t time.Time) string {
	return t.UTC().Format("20060102")
}

// tableName 返回完整表名，如 "logs_20260501"
func tableName(key string) string {
	return "logs_" + key
}

// ensureTable 确保分表及索引存在，返回对应的 INSERT prepared statement。
// 调用方须已持有 s.mu。
func (s *dbSink) ensureTable(key string) (*sql.Stmt, error) {
	if stmt, ok := s.stmtMap[key]; ok {
		return stmt, nil
	}

	tbl := tableName(key)
	_, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id         INTEGER  PRIMARY KEY AUTOINCREMENT,
			timestamp  DATETIME NOT NULL,
			level      TEXT     NOT NULL,
			caller     TEXT,
			message    TEXT,
			stacktrace TEXT,
			fields     TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_%s_ts    ON %s(timestamp);
		CREATE INDEX IF NOT EXISTS idx_%s_level ON %s(level);
	`, tbl, key, tbl, key, tbl))
	if err != nil {
		return nil, fmt.Errorf("建表 %s 失败: %w", tbl, err)
	}

	stmt, err := s.db.Prepare(fmt.Sprintf(
		`INSERT INTO %s(timestamp, level, caller, message, stacktrace, fields)
		 VALUES (?, ?, ?, ?, ?, ?)`, tbl,
	))
	if err != nil {
		return nil, fmt.Errorf("prepare %s 失败: %w", tbl, err)
	}

	s.stmtMap[key] = stmt
	return stmt, nil
}

// write 投递日志到异步队列（非阻塞）
func (s *dbSink) write(e DBEntry) {
	select {
	case s.queue <- e:
	default:
		// 队列满则静默丢弃，不影响业务
	}
}

// loop 是唯一的写入 goroutine，保证 SQLite 单写无竞争
func (s *dbSink) loop() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.cfg.FlushEvery)
	defer ticker.Stop()

	batch := make([]DBEntry, 0, s.cfg.BatchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		_ = s.writeBatch(batch)
		batch = batch[:0]
	}

	for {
		select {
		case e := <-s.queue:
			batch = append(batch, e)
			if len(batch) >= s.cfg.BatchSize {
				flush()
			}

		case <-ticker.C:
			flush()
			if s.cfg.Retention > 0 {
				s.purgeOldTables()
			}

		case <-s.done:
			// 优雅关闭：排空队列后写入剩余数据
		drain:
			for {
				select {
				case e := <-s.queue:
					batch = append(batch, e)
				default:
					break drain
				}
			}
			flush()
			return
		}
	}
}

// writeBatch 将一批日志按所属分表分组，各组各自开事务批量写入
func (s *dbSink) writeBatch(entries []DBEntry) error {
	// 按分表 key 分组，通常一批全在同一天
	groups := make(map[string][]DBEntry, 2)
	for _, e := range entries {
		key := tableKey(e.Time)
		groups[key] = append(groups[key], e)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, group := range groups {
		stmt, err := s.ensureTable(key)
		if err != nil {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			continue
		}
		txStmt := tx.Stmt(stmt)
		ok := true

		for _, e := range group {
			if _, err = txStmt.Exec(
				e.Time.UTC().Format(time.RFC3339Nano),
				e.Level,
				e.Caller,
				e.Message,
				e.Stacktrace,
				e.Fields,
			); err != nil {
				ok = false
				break
			}
		}

		if ok {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}
	return nil
}

// purgeOldTables 删除超过 Retention 天的分表（整表 DROP，效率远高于逐行 DELETE）
func (s *dbSink) purgeOldTables() {
	// cutoff 格式与表名后缀一致："YYYYMMDD"
	cutoff := time.Now().UTC().AddDate(0, 0, -s.cfg.Retention).Format("20060102")

	rows, err := s.db.Query(
		`SELECT name FROM sqlite_master
		 WHERE type='table' AND name LIKE 'logs_%'
		   AND substr(name, 6) < ?`, cutoff,
	)
	if err != nil {
		return
	}
	var tables []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			tables = append(tables, name)
		}
	}
	rows.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, tbl := range tables {
		_, _ = s.db.Exec("DROP TABLE IF EXISTS " + tbl)
		key := tbl[5:] // 去掉 "logs_" 前缀
		if stmt, ok := s.stmtMap[key]; ok {
			_ = stmt.Close()
			delete(s.stmtMap, key)
		}
	}
}

// query 跨分表查询，自动覆盖 start~end 之间的所有分表
func (s *dbSink) query(level string, start, end time.Time, limit int) ([]DBEntry, error) {
	if limit <= 0 {
		limit = 200
	}
	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = end.AddDate(0, 0, -1)
	}

	// 收集时间范围覆盖的分表名
	tableSet := make(map[string]struct{})
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		tableSet[tableName(tableKey(d))] = struct{}{}
	}

	// 只查实际存在的表
	existTables, err := s.existingTables()
	if err != nil {
		return nil, err
	}

	var results []DBEntry
	remaining := limit

	for tbl := range tableSet {
		if _, ok := existTables[tbl]; !ok {
			continue
		}

		q := fmt.Sprintf(
			`SELECT timestamp, level, caller, message, stacktrace, fields
			 FROM %s WHERE timestamp BETWEEN ? AND ?`, tbl,
		)
		args := []interface{}{
			start.UTC().Format(time.RFC3339Nano),
			end.UTC().Format(time.RFC3339Nano),
		}
		if level != "" {
			q += " AND level = ?"
			args = append(args, level)
		}
		q += " ORDER BY timestamp DESC LIMIT ?"
		args = append(args, remaining)

		rows, err := s.db.Query(q, args...)
		if err != nil {
			continue
		}
		for rows.Next() {
			var e DBEntry
			var ts string
			if err = rows.Scan(&ts, &e.Level, &e.Caller, &e.Message, &e.Stacktrace, &e.Fields); err != nil {
				continue
			}
			e.Time, _ = time.Parse(time.RFC3339Nano, ts)
			results = append(results, e)
		}
		rows.Close()

		remaining = limit - len(results)
		if remaining <= 0 {
			break
		}
	}
	return results, nil
}

// existingTables 返回数据库中实际存在的 logs_* 表集合
func (s *dbSink) existingTables() (map[string]struct{}, error) {
	rows, err := s.db.Query(
		`SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'logs_%'`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]struct{})
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			m[name] = struct{}{}
		}
	}
	return m, rows.Err()
}

func (s *dbSink) close() error {
	s.once.Do(func() { close(s.done) })
	s.wg.Wait()

	s.mu.Lock()
	for _, stmt := range s.stmtMap {
		_ = stmt.Close()
	}
	s.mu.Unlock()

	return s.db.Close()
}
