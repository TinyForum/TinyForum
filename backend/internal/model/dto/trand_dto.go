package dto

type AdminStatsTrendRequest struct {
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
	Type      string `form:"type"       binding:"required,oneof=users posts comments"`
	Interval  string `form:"interval"   binding:"omitempty,oneof=day week month"`
}
