package service

import (
	"context"
	"errors"

	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/domain"
	"github.com/codepzj/stellux/server/internal/post/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IPostService interface {
	AdminCreatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePost(ctx context.Context, post *domain.Post) error
	AdminUpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	AdminSoftDeletePost(ctx context.Context, id bson.ObjectID) error
	AdminSoftDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error
	AdminDeletePost(ctx context.Context, id bson.ObjectID) error
	AdminDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error
	AdminRestorePost(ctx context.Context, id bson.ObjectID) error
	AdminRestorePostBatch(ctx context.Context, ids []bson.ObjectID) error
	GetPostById(ctx context.Context, id bson.ObjectID) (*domain.Post, error)
	GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetPostDetailById(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error)
	GetPostDetailList(ctx context.Context, page *apiwrap.Page, postType string) ([]*domain.PostDetail, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*domain.Post, error)
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

type PostService struct {
	repo repository.IPostRepository
}

func (s *PostService) AdminCreatePost(ctx context.Context, post *domain.Post) error {
	return s.repo.Create(ctx, post)
}

func (s *PostService) AdminUpdatePost(ctx context.Context, post *domain.Post) error {
	return s.repo.Update(ctx, post)
}

func (s *PostService) AdminUpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	return s.repo.UpdatePublishStatus(ctx, id, isPublish)
}

func (s *PostService) AdminSoftDeletePost(ctx context.Context, id bson.ObjectID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *PostService) AdminSoftDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	return s.repo.SoftDeleteBatch(ctx, ids)
}

func (s *PostService) AdminDeletePost(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *PostService) AdminDeletePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	return s.repo.DeleteBatch(ctx, ids)
}

func (s *PostService) AdminRestorePost(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Restore(ctx, id)
}

func (s *PostService) AdminRestorePostBatch(ctx context.Context, ids []bson.ObjectID) error {
	return s.repo.RestoreBatch(ctx, ids)
}

func (s *PostService) GetPostById(ctx context.Context, id bson.ObjectID) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) GetPostByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	return s.repo.GetByKeyWord(ctx, keyWord)
}

func (s *PostService) GetPostDetailById(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error) {
	return s.repo.GetDetailByID(ctx, id)
}

func (s *PostService) GetPostDetailList(ctx context.Context, page *apiwrap.Page, postType string) ([]*domain.PostDetail, int64, error) {
	switch postType {
	case "publish":
		return s.repo.GetDetailList(ctx, page, "publish")
	case "draft":
		return s.repo.GetDetailList(ctx, page, "draft")
	case "bin":
		return s.repo.GetDetailList(ctx, page, "bin")
	default:
		return nil, 0, errors.New("invalid post type")
	}
}

// GetAllPublishPost 获取所有发布文章
func (s *PostService) GetAllPublishPost(ctx context.Context) ([]*domain.Post, error) {
	return s.repo.GetAllPublishPost(ctx)
}
