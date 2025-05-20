package web

import (
	"time"

	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/domain"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostVO struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	CategoryID  string    `json:"category_id"`
	TagsID      []string  `json:"tags_id"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type PostDetailVO struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	IsPublish   bool      `json:"is_publish"`
	IsTop       bool      `json:"is_top"`
	Thumbnail   string    `json:"thumbnail"`
}

type SiteMapVO struct {
	ID          string    `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetCategoryNameFromLabel(label label.Domain) string {
	return label.Name
}

func GetTagNamesFromLabels(labels []label.Domain) []string {
	return lo.Map(labels, func(label label.Domain, _ int) string {
		return label.Name
	})
}

func (h *PostHandler) PostDTOToDomain(postReq PostRequest) *domain.Post {
	return &domain.Post{
		CreatedAt:   postReq.CreatedAt,
		Title:       postReq.Title,
		Content:     postReq.Content,
		Description: postReq.Description,
		Author:      postReq.Author,
		CategoryID:  apiwrap.ConvertBsonID(postReq.CategoryID),
		TagsID:      apiwrap.ConvertBsonIDList(postReq.TagsID),
		IsPublish:   postReq.IsPublish,
		IsTop:       postReq.IsTop,
		Thumbnail:   postReq.Thumbnail,
	}
}

func (h *PostHandler) PostUpdateDTOToDomain(postUpdateReq PostUpdateRequest) *domain.Post {
	return &domain.Post{
		ID:          apiwrap.ConvertBsonID(postUpdateReq.ID),
		CreatedAt:   postUpdateReq.CreatedAt,
		Title:       postUpdateReq.Title,
		Content:     postUpdateReq.Content,
		Description: postUpdateReq.Description,
		Author:      postUpdateReq.Author,
		CategoryID:  apiwrap.ConvertBsonID(postUpdateReq.CategoryID),
		TagsID:      apiwrap.ConvertBsonIDList(postUpdateReq.TagsID),
		IsPublish:   postUpdateReq.IsPublish,
		IsTop:       postUpdateReq.IsTop,
		Thumbnail:   postUpdateReq.Thumbnail,
	}
}

func (h *PostHandler) PostDetailToVO(post *domain.PostDetail) *PostDetailVO {
	return &PostDetailVO{
		ID:          post.ID.Hex(),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Category:    GetCategoryNameFromLabel(post.Category),
		Tags:        GetTagNamesFromLabels(post.Tags),
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (h *PostHandler) PostDetailListToVOList(posts []*domain.PostDetail) []*PostDetailVO {
	return lo.Map(posts, func(post *domain.PostDetail, _ int) *PostDetailVO {
		return h.PostDetailToVO(post)
	})
}

func (h *PostHandler) PostToVO(post *domain.Post) *PostVO {
	return &PostVO{
		ID:          post.ID.Hex(),
		CreatedAt:   post.CreatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		CategoryID:  post.CategoryID.Hex(),
		TagsID: lo.Map(post.TagsID, func(id bson.ObjectID, _ int) string {
			return id.Hex()
		}),
		IsPublish: post.IsPublish,
		IsTop:     post.IsTop,
		Thumbnail: post.Thumbnail,
	}
}

func (h *PostHandler) PostListToVOList(posts []*domain.Post) []*PostVO {
	return lo.Map(posts, func(post *domain.Post, _ int) *PostVO {
		return h.PostToVO(post)
	})
}

func (h *PostHandler) PostListToSiteMapVOList(posts []*domain.Post) []*SiteMapVO {
	return lo.Map(posts, func(post *domain.Post, _ int) *SiteMapVO {
		return &SiteMapVO{
			ID:          post.ID.Hex(),
			UpdatedAt:   post.UpdatedAt,
		}
	})
}
