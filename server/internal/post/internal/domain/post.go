package domain

import (
	"time"

	"github.com/codepzj/stellux/server/internal/label"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Post 文章
type Post struct {
	ID          bson.ObjectID   // 文章ID
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	Title       string          // 标题
	Content     string          // 内容
	Description string          // 描述
	Author      string          // 作者
	CategoryID  bson.ObjectID   // 分类ID
	TagsID      []bson.ObjectID // 标签ID
	IsPublish   bool            // 是否发布
	IsTop       bool            // 是否置顶
	Thumbnail   string          // 缩略图
}

type PostDetail struct {
	ID          bson.ObjectID  // 文章ID
	CreatedAt   time.Time      // 创建时间
	UpdatedAt   time.Time      // 更新时间
	Title       string         // 标题
	Content     string         // 内容
	Description string         // 描述
	Author      string         // 作者
	Category    label.Domain   // 分类
	Tags        []label.Domain // 标签
	IsPublish   bool           // 是否发布
	IsTop       bool           // 是否置顶
	Thumbnail   string         // 缩略图
}

// Page 分页
type Page struct {
	PageNo   int64  // 页码
	PageSize int64  // 每页条数
	Keyword  string // 关键词
	Field    string // 排序字段
	Order    string // 排序方向
}
