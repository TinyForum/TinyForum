package vo

type AdminSetUserRole struct {
	Message    string `json:"message"`
	UserID     uint64 `json:"user_id"` // 根据实际类型调整
	NewRole    string `json:"new_role"`
	OperatorID uint   `json:"operator_id"` // 根据实际类型调整
}
