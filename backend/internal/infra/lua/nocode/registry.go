// // nocode/registry.go
// package nocode

// import (
// 	"sort"
// 	"sync"
// )

// type NodeRegistry struct {
// 	mu         sync.RWMutex
// 	triggers   map[string]NodeMeta
// 	conditions map[string]NodeMeta
// 	actions    map[string]NodeMeta
// }

// var globalRegistry = &NodeRegistry{
// 	triggers:   make(map[string]NodeMeta),
// 	conditions: make(map[string]NodeMeta),
// 	actions:    make(map[string]NodeMeta),
// }

// // RegisterTrigger 注册一个触发器节点
// func RegisterTrigger(meta NodeMeta) {
// 	globalRegistry.mu.Lock()
// 	defer globalRegistry.mu.Unlock()
// 	globalRegistry.triggers[meta.Type] = meta
// }

// // RegisterCondition 注册一个条件节点
// func RegisterCondition(meta NodeMeta) {
// 	globalRegistry.mu.Lock()
// 	defer globalRegistry.mu.Unlock()
// 	globalRegistry.conditions[meta.Type] = meta
// }

// // RegisterAction 注册一个动作节点
// func RegisterAction(meta NodeMeta) {
// 	globalRegistry.mu.Lock()
// 	defer globalRegistry.mu.Unlock()
// 	globalRegistry.actions[meta.Type] = meta
// }

// // GetNocodeMetadata 收集所有已注册节点，返回给前端
// func GetNocodeMetadata() *NocodeMetadata {
// 	globalRegistry.mu.RLock()
// 	defer globalRegistry.mu.RUnlock()

// 	triggers := make([]NodeMeta, 0, len(globalRegistry.triggers))
// 	for _, v := range globalRegistry.triggers {
// 		triggers = append(triggers, v)
// 	}
// 	conditions := make([]NodeMeta, 0, len(globalRegistry.conditions))
// 	for _, v := range globalRegistry.conditions {
// 		conditions = append(conditions, v)
// 	}
// 	actions := make([]NodeMeta, 0, len(globalRegistry.actions))
// 	for _, v := range globalRegistry.actions {
// 		actions = append(actions, v)
// 	}

// 	// 按名称排序，保证输出稳定
// 	sort.Slice(triggers, func(i, j int) bool { return triggers[i].Name < triggers[j].Name })
// 	sort.Slice(conditions, func(i, j int) bool { return conditions[i].Name < conditions[j].Name })
// 	sort.Slice(actions, func(i, j int) bool { return actions[i].Name < actions[j].Name })

// 	return &NocodeMetadata{
// 		Triggers:   triggers,
// 		Conditions: conditions,
// 		Actions:    actions,
// 	}
// }

package nocode
