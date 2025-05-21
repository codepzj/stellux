package service

import (
	"context"

	"github.com/codepzj/stellux/server/internal/document/internal/domain"
	"github.com/codepzj/stellux/server/internal/document/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDocumentService interface {
	Create(ctx context.Context, doc *domain.Document) error
	FindAllRoot(ctx context.Context) ([]*domain.Document, error)
	FindAllParent(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error)
	UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
	DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
	RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
}

var _ IDocumentService = (*DocumentService)(nil)

func NewDocumentService(repo repository.IDocumentRepository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

type DocumentService struct {
	repo repository.IDocumentRepository
}

func (s *DocumentService) Create(ctx context.Context, doc *domain.Document) error {
	return s.repo.Create(ctx, doc)
}

func (s *DocumentService) FindAllRoot(ctx context.Context) ([]*domain.Document, error) {
	return s.repo.FindAllByType(ctx, "root")
}

func (s *DocumentService) FindAllParent(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	return s.repo.FindAllByTypeAndDocumentID(ctx, "parent", documentID)
}

func (s *DocumentService) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	return s.repo.FindAllByDocumentID(ctx, documentID)
}

func (s *DocumentService) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*domain.Document, error) {
	return s.repo.GetDocumentByID(ctx, id)
}

func (s *DocumentService) UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	return s.repo.UpdateDocumentByID(ctx, id, title, content)
}

func (s *DocumentService) RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	return s.repo.RenameDocumentByID(ctx, id, title)
}

func (s *DocumentService) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *DocumentService) DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	return s.repo.DeleteByIDList(ctx, idList)
}