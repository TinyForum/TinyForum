package request

type ResolveTaskInput struct {
	Note string `json:"note" binding:"max=500"`
}
