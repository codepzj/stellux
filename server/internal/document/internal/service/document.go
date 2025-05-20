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
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error)
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
	return s.repo.FindAllRoot(ctx)
}

func (s *DocumentService) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*domain.Document, error) {
	return s.repo.FindAllByDocumentID(ctx, documentID)
}
