import { useState, useEffect } from "react";

// ==================== 插件日志 ====================
export function PluginLogsTab() {
  const [logs, setLogs] = useState<
    { time: string; level: string; plugin: string; message: string }[]
  >([]);
  // 模拟获取日志（后续对接真实 API）
  useEffect(() => {
    setLogs([
      {
        time: "2025-01-15 10:32:21",
        level: "info",
        plugin: "SEO Enhancer",
        message: "初始化成功",
      },
      {
        time: "2025-01-15 10:28:15",
        level: "error",
        plugin: "Analytics",
        message: "网络请求失败",
      },
    ]);
  }, []);
  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-base font-semibold">插件运行日志</h3>
        <button className="btn btn-xs btn-ghost">清空日志</button>
      </div>
      <div className="overflow-x-auto">
        <table className="table table-xs">
          <thead>
            <tr>
              <th>时间</th>
              <th>级别</th>
              <th>插件</th>
              <th>消息</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log, idx) => (
              <tr key={idx}>
                <td>{log.time}</td>
                <td>
                  <span
                    className={`badge badge-sm ${log.level === "error" ? "badge-error" : "badge-info"}`}
                  >
                    {log.level}
                  </span>
                </td>
                <td>{log.plugin}</td>
                <td>{log.message}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
