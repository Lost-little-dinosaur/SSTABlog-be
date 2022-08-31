package draftBox

type AddDraftRequest struct {
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type UpdateDraftRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type GetAllDraftResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
}
