package request

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
