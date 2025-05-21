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
		adminGroup.GET("/:id", apiwrap.Wrap(h.GetDocumentByID))
		adminGroup.GET("/list", apiwrap.Wrap(h.GetAllRootDoc))
		adminGroup.GET("/parent-list", apiwrap.Wrap(h.GetAllParentDoc))
		adminGroup.POST("/create", apiwrap.WrapWithBody(h.CreateDocument))
		adminGroup.PUT("/save", apiwrap.WrapWithBody(h.SaveDocument))
		adminGroup.PUT("/rename", apiwrap.WrapWithBody(h.RenameDocument))
		adminGroup.DELETE("/delete", apiwrap.WrapWithBody(h.DeleteDocument))
		adminGroup.DELETE("/delete-list", apiwrap.WrapWithBody(h.DeleteDocumentList))
	}
}

func (h *DocumentHandler) CreateDocument(c *gin.Context, documentReq DocumentRequest) *apiwrap.Response[any] {
	err := h.serv.Create(c.Request.Context(), h.DocumentRequestToDomain(documentReq))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

func (h *DocumentHandler) SaveDocument(c *gin.Context, req UpdateDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.UpdateDocumentByID(c.Request.Context(), apiwrap.ConvertBsonID(req.ID), req.Title, req.Content)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

func (h *DocumentHandler) RenameDocument(c *gin.Context, req RenameDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.RenameDocumentByID(c.Request.Context(), apiwrap.ConvertBsonID(req.ID), req.Title)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context, req DeleteDocumentRequest) *apiwrap.Response[any] {
	err := h.serv.DeleteByID(c.Request.Context(), apiwrap.ConvertBsonID(req.DocumentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.Success()
}

func (h *DocumentHandler) DeleteDocumentList(c *gin.Context, req DeleteDocumentListRequest) *apiwrap.Response[any] {
	err := h.serv.DeleteByIDList(c.Request.Context(), apiwrap.ConvertBsonIDList(req.DocumentIDList))
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

func (h *DocumentHandler) GetDocumentByID(c *gin.Context) *apiwrap.Response[any] {
	documentID := c.Param("id")
	document, err := h.serv.GetDocumentByID(c.Request.Context(), apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainToVO(document), "获取文档成功")
}

func (h *DocumentHandler) GetAllRootDoc(c *gin.Context) *apiwrap.Response[any] {
	documentList, err := h.serv.FindAllRoot(c.Request.Context())
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToRootVOList(documentList), "获取文档列表成功")
}

func (h *DocumentHandler) GetAllParentDoc(c *gin.Context) *apiwrap.Response[any] {
	documentID := c.Query("document_id")
	documentList, err := h.serv.FindAllParent(c.Request.Context(), apiwrap.ConvertBsonID(documentID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.DocumentDomainListToRootVOList(documentList), "获取文档列表成功")
}

func (h *DocumentHandler) DocumentRequestToDomain(req DocumentRequest) *domain.Document {
	var parentID, documentID bson.ObjectID
	if req.DocumentType == "root" {
		// 根节点不需要parent_id和document_id
		parentID = bson.ObjectID{}
		documentID = bson.ObjectID{}
	} else {
		parentID = apiwrap.ConvertBsonID(req.ParentID)
		documentID = apiwrap.ConvertBsonID(req.DocumentID)
	}
	return &domain.Document{
		Title:        req.Title,
		Content:      req.Content,
		Description:  req.Description,
		Thumbnail:    req.Thumbnail,
		DocumentType: req.DocumentType,
		IsPublic:     req.IsPublic,
		ParentID:     parentID,
		DocumentID:   documentID,
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

func (h *DocumentHandler) DocumentDomainToRootVO(doc *domain.Document) *DocumentRootVO {
	return &DocumentRootVO{
		ID:          doc.ID.Hex(),
		CreatedAt:   doc.CreatedAt.String(),
		UpdatedAt:   doc.UpdatedAt.String(),
		Title:       doc.Title,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		IsPublic:    doc.IsPublic,
	}
}

func (h *DocumentHandler) DocumentDomainListToRootVOList(docList []*domain.Document) []*DocumentRootVO {
	return lo.Map(docList, func(doc *domain.Document, _ int) *DocumentRootVO {
		return h.DocumentDomainToRootVO(doc)
	})
}

func (h *DocumentHandler) DocumentDomainToVO(doc *domain.Document) *DocumentVO {
	return &DocumentVO{
		ID: doc.ID.Hex(),
		CreatedAt: doc.CreatedAt.String(),
		UpdatedAt: doc.UpdatedAt.String(),
		Title: doc.Title,
		Content: doc.Content,
		DocumentType: doc.DocumentType,
		IsPublic: doc.IsPublic,
		ParentID: apiwrap.BsonID(doc.ParentID.Hex()),
		DocumentID: apiwrap.BsonID(doc.DocumentID.Hex()),
	}
}

func (h *DocumentHandler) DocumentDomainListToVOList(docList []*domain.Document) []*DocumentVO {
	return lo.Map(docList, func(doc *domain.Document, _ int) *DocumentVO {
		return h.DocumentDomainToVO(doc)
	})
}
