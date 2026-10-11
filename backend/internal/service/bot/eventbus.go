package bot

import (
	"sync"

	"tiny-forum/pkg/logger"
)

// ─── EventBus ─────────────────────────────────────────────────────────────

// EventBus 进程内事件总线。subscribers 会被并发读写，必须用 RWMutex 保护。
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]func(map[string]any)
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]func(map[string]any))}
}

func (e *EventBus) Subscribe(event string, fn func(map[string]any)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.subscribers[event] = append(e.subscribers[event], fn)
}

func (e *EventBus) Publish(event string, data map[string]any) {
	// 复制订阅者切片后再执行，避免持锁运行回调造成死锁/阻塞
	e.mu.RLock()
	subs := make([]func(map[string]any), len(e.subscribers[event]))
	copy(subs, e.subscribers[event])
	e.mu.RUnlock()

	for _, fn := range subs {
		go func(fn func(map[string]any)) {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("bot event handler panic: %v", r)
				}
			}()
			fn(data)
		}(fn)
	}
}

func (s *service) PublishEvent(eventName string, data map[string]any) {
	s.eventBus.Publish(eventName, data)
}
