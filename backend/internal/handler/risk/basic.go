package risk

import (
	"tiny-forum/internal/service/check"
	riskservice "tiny-forum/internal/service/risk"
)

type RiskHandler struct {
	checkSvc check.ContentCheckService
	riskSvc  riskservice.RiskService
}

func NewRiskHandler(checkSvc check.ContentCheckService, riskSvc riskservice.RiskService) *RiskHandler {
	return &RiskHandler{checkSvc: checkSvc, riskSvc: riskSvc}
}
