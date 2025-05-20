package repository

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/bsonx"
	"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/domain"
	"github.com/codepzj/stellux/server/internal/post/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IPostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	UpdatePublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	SoftDelete(ctx context.Context, id bson.ObjectID) error
	SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Restore(ctx context.Context, id bson.ObjectID) error
	RestoreBatch(ctx context.Context, ids []bson.ObjectID) error
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.Post, error)
	GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetDetailByID(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error)
	GetDetailList(ctx context.Context, page *apiwrap.Page, postType string) ([]*domain.PostDetail, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*domain.Post, error)
}

var _ IPostRepository = (*PostRepository)(nil)

func NewPostRepository(dao dao.IPostDao) *PostRepository {
	return &PostRepository{dao: dao}
}

type PostRepository struct {
	dao dao.IPostDao
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.dao.Create(ctx, r.PostDomainToPostDO(post))
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.dao.Update(ctx, post.ID, r.PostDomainToUpdatePostDO(post))
}

// UpdatePostPublishStatus 更新文章发布状态
func (r *PostRepository) UpdatePublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	return r.dao.UpdatePostPublishStatus(ctx, id, isPublish)
}

func (r *PostRepository) SoftDelete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.SoftDelete(ctx, id)
}

// SoftDeleteBatch 批量软删除文章
func (r *PostRepository) SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.SoftDeleteBatch(ctx, ids)
}

// Delete 删除文章
func (r *PostRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

// DeleteBatch 批量删除文章
func (r *PostRepository) DeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.DeleteBatch(ctx, ids)
}

// Restore 恢复文章
func (r *PostRepository) Restore(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Restore(ctx, id)
}

// RestoreBatch 批量恢复文章
func (r *PostRepository) RestoreBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.RestoreBatch(ctx, ids)
}

// GetByID 获取文章
func (r *PostRepository) GetByID(ctx context.Context, id bson.ObjectID) (*domain.Post, error) {
	post, err := r.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomain(post), nil
}

// GetDetailByID 获取文章
func (r *PostRepository) GetDetailByID(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error) {
	postCategoryTags, err := r.dao.GetDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostCategoryTagsDOToPostDetail(postCategoryTags), nil
}

// GetByKeyWord 获取文章
func (r *PostRepository) GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	posts, err := r.dao.GetByKeyWord(ctx, keyWord)
	if err != nil {
		return nil, err
	}
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	}), nil
}

// GetDetailList 获取文章列表
func (r *PostRepository) GetDetailList(ctx context.Context, page *apiwrap.Page, postType string) ([]*domain.PostDetail, int64, error) {
	var cond bson.D
	switch postType {
	case "publish":
		cond = query.NewBuilder().Or(query.RegexOptions("title", page.Keyword, "i"), query.RegexOptions("description", page.Keyword, "i")).And(query.Eq("deleted_at", nil), query.Eq("is_publish", true)).Build()
	case "draft":
		cond = query.NewBuilder().Or(query.RegexOptions("title", page.Keyword, "i"), query.RegexOptions("description", page.Keyword, "i")).And(query.Eq("deleted_at", nil), query.Eq("is_publish", false)).Build()
	case "bin":
		cond = query.NewBuilder().Or(query.RegexOptions("title", page.Keyword, "i"), query.RegexOptions("description", page.Keyword, "i")).And(query.Ne("deleted_at", nil)).Build()
	}
	skip := (page.PageNo - 1) * page.PageSize
	limit := page.PageSize
	sortBuilder := bsonx.NewD().Add("is_top", -1).Add("created_at", -1)
	if page.Field != "" {
		sortBuilder.Add(page.Field, r.OrderConvertToInt(page.Order))
	}
	pagePipeline := aggregation.NewStageBuilder().Match(cond).Lookup("label", "category", &aggregation.LookUpOptions{
		LocalField:   "category_id",
		ForeignField: "_id",
	}).Unwind("$category", nil).Lookup("label", "tags", &aggregation.LookUpOptions{
		LocalField:   "tags_id",
		ForeignField: "_id",
	}).Skip(skip).Limit(limit).Sort(sortBuilder.Build()).Build()

	posts, count, err := r.dao.GetDetailList(ctx, pagePipeline, cond)
	if err != nil {
		return nil, 0, err
	}
	return r.PostCategoryTagsDOToPostDetailList(posts), count, nil
}

func (r *PostRepository) GetAllPublishPost(ctx context.Context) ([]*domain.Post, error) {
	posts, err := r.dao.GetAllPublishPost(ctx)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomainList(posts), nil
}

// PostDomain2PostDO 将domain.Post转换为dao.Post
func (r *PostRepository) PostDomainToPostDO(post *domain.Post) *dao.Post {
	return &dao.Post{
		Model:       mongox.Model{CreatedAt: post.CreatedAt},
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		CategoryID:  post.CategoryID,
		TagsID:      post.TagsID,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDomainToUpdatePostDO(post *domain.Post) *dao.UpdatePost {
	return &dao.UpdatePost{
		CreatedAt:   post.CreatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		CategoryID:  post.CategoryID,
		TagsID:      post.TagsID,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDOToPostDomainList(posts []*dao.Post) []*domain.Post {
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	})
}

// PostDOToPostDomain 将dao.Post转换为domain.Post
func (r *PostRepository) PostDOToPostDomain(post *dao.Post) *domain.Post {
	return &domain.Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		CategoryID:  post.CategoryID,
		TagsID:      post.TagsID,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

// PostCategoryTagsDOToPostDetail 将dao.PostCategoryTags转换为domain.PostDetail
func (r *PostRepository) PostCategoryTagsDOToPostDetail(postCategoryTags *dao.PostCategoryTags) *domain.PostDetail {
	return &domain.PostDetail{
		ID:          postCategoryTags.ID,
		CreatedAt:   postCategoryTags.CreatedAt,
		UpdatedAt:   postCategoryTags.UpdatedAt,
		Title:       postCategoryTags.Title,
		Content:     postCategoryTags.Content,
		Description: postCategoryTags.Description,
		Author:      postCategoryTags.Author,
		Category:    postCategoryTags.Category,
		Tags:        postCategoryTags.Tags,
		Thumbnail:   postCategoryTags.Thumbnail,
		IsPublish:   postCategoryTags.IsPublish,
		IsTop:       postCategoryTags.IsTop,
	}
}

func (r *PostRepository) PostCategoryTagsDOToPostDetailList(postCategoryTags []*dao.PostCategoryTags) []*domain.PostDetail {
	return lo.Map(postCategoryTags, func(postCategoryTag *dao.PostCategoryTags, _ int) *domain.PostDetail {
		return r.PostCategoryTagsDOToPostDetail(postCategoryTag)
	})
}

func (r *PostRepository) OrderConvertToInt(order string) int {
	switch order {
	case "DESC":
		return -1
	case "ASC":
		return 1
	default:
		return -1
	}
}
