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

type DeleteDocumentRequest struct {
	DocumentID string `json:"document_id" binding:"required"`
}

type DeleteDocumentListRequest struct {
	DocumentIDList []string `json:"document_id_list" binding:"required"`
}

type UpdateDocumentRequest struct {
	ID       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type RenameDocumentRequest struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}
