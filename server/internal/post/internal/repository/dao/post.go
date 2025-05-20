package dao

import (
	"context"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/codepzj/stellux/server/internal/label"
	"github.com/codepzj/stellux/server/internal/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Post struct {
	mongox.Model `bson:",inline"`
	Title        string          `bson:"title"`
	Content      string          `bson:"content"`
	Description  string          `bson:"description"`
	Author       string          `bson:"author"`
	CategoryID   bson.ObjectID   `bson:"category_id"`
	TagsID       []bson.ObjectID `bson:"tags_id"`
	IsPublish    bool            `bson:"is_publish"`
	IsTop        bool            `bson:"is_top"`
	Thumbnail    string          `bson:"thumbnail"`
}

type PostUpdate struct {
	Title       string          `bson:"title"`
	Content     string          `bson:"content"`
	Description string          `bson:"description"`
	Author      string          `bson:"author"`
	CategoryID  bson.ObjectID   `bson:"category_id"`
	TagsID      []bson.ObjectID `bson:"tags_id"`
	IsPublish   bool            `bson:"is_publish"`
	IsTop       bool            `bson:"is_top"`
	Thumbnail   string          `bson:"thumbnail"`
}

// 聚合查询返回带有category和tags的结构体
type PostCategoryTags struct {
	ID          bson.ObjectID  `bson:"_id"`
	CreatedAt   time.Time      `bson:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at"`
	Title       string         `bson:"title"`
	Content     string         `bson:"content"`
	Description string         `bson:"description"`
	Author      string         `bson:"author"`
	Category    label.Domain   `bson:"category"`
	Tags        []label.Domain `bson:"tags"`
	IsPublish   bool           `bson:"is_publish"`
	IsTop       bool           `bson:"is_top"`
	Thumbnail   string         `bson:"thumbnail"`
}

type UpdatePost struct {
	CreatedAt   time.Time       `bson:"created_at,omitempty"`
	Title       string          `bson:"title"`
	Content     string          `bson:"content"`
	Description string          `bson:"description"`
	Author      string          `bson:"author"`
	CategoryID  bson.ObjectID   `bson:"category_id"`
	TagsID      []bson.ObjectID `bson:"tags_id"`
	IsPublish   bool            `bson:"is_publish"`
	IsTop       bool            `bson:"is_top"`
	Thumbnail   string          `bson:"thumbnail"`
}

type IPostDao interface {
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, id bson.ObjectID, post *UpdatePost) error
	UpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	SoftDelete(ctx context.Context, id bson.ObjectID) error
	SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Restore(ctx context.Context, id bson.ObjectID) error
	RestoreBatch(ctx context.Context, ids []bson.ObjectID) error
	GetByID(ctx context.Context, id bson.ObjectID) (*Post, error)
	GetByKeyWord(ctx context.Context, keyWord string) ([]*Post, error)
	GetDetailByID(ctx context.Context, id bson.ObjectID) (*PostCategoryTags, error)
	GetDetailList(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D) ([]*PostCategoryTags, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*Post, error)
}

var _ IPostDao = (*PostDao)(nil)

func NewPostDao(db *mongox.Database) *PostDao {
	return &PostDao{coll: mongox.NewCollection[Post](db, "post")}
}

type PostDao struct {
	coll *mongox.Collection[Post]
}

// Create 创建文章
func (d *PostDao) Create(ctx context.Context, post *Post) error {
	insertResult, err := d.coll.Creator().InsertOne(ctx, post)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.Wrap(err, "插入文章失败")
	}
	return nil
}

func (d *PostDao) Update(ctx context.Context, id bson.ObjectID, post *UpdatePost) error {
	updateResult, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().SetFields(post).Build()).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("文章修改失败")
	}
	return nil
}

// UpdatePostPublishStatus 更新文章发布状态
func (d *PostDao) UpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("is_publish", isPublish).Build()).UpdateOne(ctx)
	return err
}

// SoftDelete 软删除文章
func (d *PostDao) SoftDelete(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("deleted_at", time.Now()).Set("is_publish", false).Set("is_top", false).Build()).UpdateOne(ctx)
	return err
}

// Delete 删除文章
func (d *PostDao) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	return err
}

// DeleteBatch 批量删除文章
func (d *PostDao) DeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.In("_id", ids...)).DeleteMany(ctx)
	return err
}

// Restore 恢复文章
func (d *PostDao) Restore(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("deleted_at", nil).Build()).UpdateOne(ctx)
	return err
}

// GetDetailByID 获取文章
func (d *PostDao) GetDetailByID(ctx context.Context, id bson.ObjectID) (*PostCategoryTags, error) {
	// 设置管道,聚合查询包含详细分类和标签的文章
	pipeline := aggregation.NewStageBuilder().Match(query.Id(id)).Lookup("label", "category", &aggregation.LookUpOptions{
		LocalField:   "category_id",
		ForeignField: "_id",
	}).Unwind("$category", nil).Lookup("label", "tags", &aggregation.LookUpOptions{
		LocalField:   "tags_id",
		ForeignField: "_id",
	}).Build()
	var postResult []PostCategoryTags
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &postResult)
	if err != nil {
		return nil, err
	}
	if len(postResult) == 0 {
		return nil, errors.New("文章不存在")
	}
	return &postResult[0], err
}

// GetByID 获取文章
func (d *PostDao) GetByID(ctx context.Context, id bson.ObjectID) (*Post, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

// GetByKeyWord 获取文章
func (d *PostDao) GetByKeyWord(ctx context.Context, keyWord string) ([]*Post, error) {
	cond := query.NewBuilder().Or(query.RegexOptions("title", keyWord, "i"), query.RegexOptions("description", keyWord, "i")).And(query.Eq("deleted_at", nil), query.Eq("is_publish", true)).Build()
	return d.coll.Finder().Filter(cond).Find(ctx)
}

// GetDetailList 获取文章列表
func (d *PostDao) GetDetailList(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D) ([]*PostCategoryTags, int64, error) {
	var postResult []PostCategoryTags
	err := d.coll.Aggregator().Pipeline(pagePipeline).AggregateWithParse(ctx, &postResult)
	if err != nil {
		return nil, 0, err
	}
	count, err := d.coll.Finder().Filter(cond).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return utils.ValToPtrList(postResult), count, err
}

// SoftDeleteBatch 批量软删除文章
func (d *PostDao) SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	_, err := d.coll.Updater().Filter(query.In("_id", ids...)).Updates(update.NewBuilder().Set("deleted_at", time.Now()).Set("is_publish", false).Set("is_top", false).Build()).UpdateMany(ctx)
	return err
}

// RestoreBatch 批量恢复文章
func (d *PostDao) RestoreBatch(ctx context.Context, ids []bson.ObjectID) error {
	_, err := d.coll.Updater().Filter(query.In("_id", ids...)).Updates(update.NewBuilder().Set("deleted_at", nil).Build()).UpdateMany(ctx)
	return err
}

// GetAllPublishPost 获取所有发布文章
func (d *PostDao) GetAllPublishPost(ctx context.Context) ([]*Post, error) {
	return d.coll.Finder().Filter(query.Eq("is_publish", true)).Find(ctx)
}
