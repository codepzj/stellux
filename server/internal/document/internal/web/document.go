package web

import (
	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func NewDocumentHandler(serv service.IDocumentService) *DocumentHandler {
	return &DocumentHandler{
		serv: serv,
	}
}

type DocumentHandler struct {
	serv service.IDocumentService
}

func (h *DocumentHandler) RegisterGinRoutes(engine *gin.Engine) {
	documentGroup := engine.Group("/document")
	{
		documentGroup.GET("/tree", apiwrap.Wrap(h.GetDocumentTreeByID))
	}
	adminGroup := engine.Group("/admin-api/document")
	{
		adminGroup.GET("/list",apiwrap.Wrap(h.GetAllRootDoc))
		adminGroup.POST("/create", apiwrap.WrapWithBody(h.CreateDocument))
	}
}

func (h *DocumentHandler) CreateDocument(c *gin.Context, documentReq DocumentRequest) *apiwrap.Response[any] {
	err := h.serv.Create(c.Request.Context(), h.DocumentRequestToDomain(documentReq))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

func (h *DocumentHandler) GetDocumentTreeByID(c *gin.Context) *apiwrap.Response[any] {
	documentID := c.Query("document_id")
	documentList, err := h.serv.FindAllByDocumentID(c.Request.Context(), apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToTreeVOList(documentList), "获取文档目录树成功")
}

func (h *DocumentHandler) GetAllRootDoc(c *gin.Context) *apiwrap.Response[any] {	
	documentList, err := h.serv.FindAllRoot(c.Request.Context())
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToTreeVOList(documentList), "获取文档列表成功")
}

func (h *DocumentHandler) DocumentRequestToDomain(req DocumentRequest) *domain.Document {
	if req.DocumentType == "root" {
		// 根节点不需要parent_id和document_id
		return &domain.Document{
			Title:        req.Title,
			Content:      req.Content,
			DocumentType: req.DocumentType,
			IsPublic:     req.IsPublic,
			ParentID:     bson.ObjectID{},
			DocumentID:   bson.ObjectID{},
		}
	}
	return &domain.Document{
		Title:        req.Title,
		Content:      req.Content,
		DocumentType: req.DocumentType,
		IsPublic:     req.IsPublic,
		ParentID:     apiwrap.ConvertBsonID(req.ParentID),
		DocumentID:   apiwrap.ConvertBsonID(req.DocumentID),
	}
}

func (h *DocumentHandler) DocumentDomainToTreeVO(doc *domain.Document) *DocumentTreeVO {
	return &DocumentTreeVO{
		ID:           doc.ID.Hex(),
		CreatedAt:    doc.CreatedAt.String(),
		UpdatedAt:    doc.UpdatedAt.String(),
		Title:        doc.Title,
		DocumentType: doc.DocumentType,
		IsPublic:     doc.IsPublic,
		ParentID:     apiwrap.BsonID(doc.ParentID.Hex()),
		DocumentID:   apiwrap.BsonID(doc.DocumentID.Hex()),
	}
}

func (h *DocumentHandler) DocumentDomainListToTreeVOList(docList []*domain.Document) []*DocumentTreeVO {
	return lo.Map(docList, func(doc *domain.Document, _ int) *DocumentTreeVO {
		return h.DocumentDomainToTreeVO(doc)
	})
}
