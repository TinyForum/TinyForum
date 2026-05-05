package plugin

import "tiny-forum/internal/service/plugin"

type PluginHandler struct {
	service plugin.PluginService
}

func NewPluginHandler(svc plugin.PluginService) *PluginHandler {
	return &PluginHandler{service: svc}
}
