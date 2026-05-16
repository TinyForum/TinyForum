package bot

// ─── EventBus ─────────────────────────────────────────────────────────────

type EventBus struct {
	subscribers map[string][]func(map[string]any)
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]func(map[string]any))}
}

func (e *EventBus) Subscribe(event string, fn func(map[string]any)) {
	e.subscribers[event] = append(e.subscribers[event], fn)
}

func (e *EventBus) Publish(event string, data map[string]any) {
	for _, fn := range e.subscribers[event] {
		go fn(data)
	}
}

func (s *service) PublishEvent(eventName string, data map[string]any) {
	s.eventBus.Publish(eventName, data)
}
