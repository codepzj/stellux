package repository

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	FindAllByType(ctx context.Context,document_type string) ([]*domain.Document, error)
    FindAllByTypeAndDocumentID(ctx context.Context,document_type string,documentID bson.ObjectID) ([]*domain.Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error)	
	DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
	RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
}

var _ IDocumentRepository = (*DocumentRepository)(nil)

func NewDocumentRepository(dao dao.IDocumentDao) *DocumentRepository {
	return &DocumentRepository{dao: dao}
}

type DocumentRepository struct {
	dao dao.IDocumentDao
}

func (r *DocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	return r.dao.Create(ctx, r.DomainToDao(doc))
}

func (r *DocumentRepository) FindAllByType(ctx context.Context, document_type string) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByType(ctx, document_type)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

func (r *DocumentRepository) FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByTypeAndDocumentID(ctx, document_type, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

func (r *DocumentRepository) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomainList(documentList), nil
}

func (r *DocumentRepository) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	document, err := r.dao.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.DaoToDomain(document), nil
}

func (r *DocumentRepository) UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	return r.dao.UpdateDocumentByID(ctx, id, title, content)
}

func (r *DocumentRepository) RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	return r.dao.RenameDocumentByID(ctx, id, title)
}

func (r *DocumentRepository) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	return r.dao.DeleteByID(ctx, id)
}

func (r *DocumentRepository) DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	return r.dao.DeleteByIDList(ctx, idList)
}

func (r *DocumentRepository) DomainToDao(doc *domain.Document) *dao.Document {
	return &dao.Document{
		Title:      doc.Title,
		Content:    doc.Content,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		DocumentType: doc.DocumentType,
		IsPublic:     doc.IsPublic,
		ParentID:   doc.ParentID,
		DocumentID: doc.DocumentID,
	}
}

func (r *DocumentRepository) DaoToDomain(doc *dao.Document) *domain.Document {
	return &domain.Document{
		ID:         doc.ID,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
		Title:      doc.Title,
		Content:    doc.Content,
		Description: doc.Description,
		Thumbnail:   doc.Thumbnail,
		DocumentType: doc.DocumentType,
		IsPublic:     doc.IsPublic,
		ParentID:   doc.ParentID,
		DocumentID: doc.DocumentID,
	}
}

func (r *DocumentRepository) DaoToDomainList(docList []*dao.Document) []*domain.Document {
	return lo.Map(docList, func(doc *dao.Document, _ int) *domain.Document {
		return r.DaoToDomain(doc)
	})
}
