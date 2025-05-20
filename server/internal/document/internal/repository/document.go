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
	FindAllRoot(ctx context.Context) ([]*domain.Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
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

func (r *DocumentRepository) FindAllRoot(ctx context.Context) ([]*domain.Document, error) {
	documentList, err := r.dao.FindAllRoot(ctx)
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

func (r *DocumentRepository) DomainToDao(doc *domain.Document) *dao.Document {
	return &dao.Document{
		Title:      doc.Title,
		Content:    doc.Content,
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
