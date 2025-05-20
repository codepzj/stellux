package web

type DocumentRequest struct {
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	DocumentType string `json:"document_type" binding:"required"`
	IsPublic     bool   `json:"is_public"`
	ParentID     string `json:"parent_id"`
	DocumentID   string `json:"document_id"`
}
