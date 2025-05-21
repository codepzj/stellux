package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	mongox.Model `bson:",inline"`
	Title        string
	Content      string
	Description  string
	Thumbnail    string
	DocumentType string        `bson:"document_type"`
	IsPublic     bool          `bson:"is_public"`
	ParentID     bson.ObjectID `bson:"parent_id,omitempty"`   // 根节点不需要parent_id
	DocumentID   bson.ObjectID `bson:"document_id,omitempty"` // 根节点不需要document_id
}

type IDocumentDao interface {
	Create(ctx context.Context, doc *Document) error
	FindAllByType(ctx context.Context, document_type string) ([]*Document, error)
	FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error)
	GetDocumentByID(ctx context.Context, id bson.ObjectID) (*Document, error)
	UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error
	RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error
	DeleteByID(ctx context.Context, id bson.ObjectID) error
	DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error
}

var _ IDocumentDao = (*DocumentDao)(nil)

func NewDocumentDao(db *mongox.Database) *DocumentDao {
	return &DocumentDao{coll: mongox.NewCollection[Document](db, "document")}
}

type DocumentDao struct {
	coll *mongox.Collection[Document]
}

func (d *DocumentDao) Create(ctx context.Context, doc *Document) error {
	if doc.DocumentType == "root" {
		insertResult, err := d.coll.Creator().InsertOne(ctx, &Document{
			Title:        doc.Title,
			Content:      doc.Content,
			Description:  doc.Description,
			Thumbnail:    doc.Thumbnail,
			DocumentType: doc.DocumentType,
			IsPublic:     doc.IsPublic,
		})
		if err != nil {
			return err
		}
		if insertResult.InsertedID == nil {
			return errors.New("新增根文档失败")
		}
		return nil
	}
	insertResult, err := d.coll.Creator().InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.New("新增文档失败")
	}
	return nil
}

func (d *DocumentDao) FindAllByType(ctx context.Context, document_type string) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", document_type).Build()).Sort(bson.M{"updated_at": -1}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

func (d *DocumentDao) FindAllByTypeAndDocumentID(ctx context.Context, document_type string, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", document_type).Eq("document_id", documentID).Build()).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

func (d *DocumentDao) FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Or(query.Id(documentID), query.Eq("document_id", documentID)).Build()).Find(ctx)
	if err != nil {
		return nil, err
	}
	return documentList, nil
}

func (d *DocumentDao) GetDocumentByID(ctx context.Context, id bson.ObjectID) (*Document, error) {
	document, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, err
	}
	return document, nil
}

func (d *DocumentDao) UpdateDocumentByID(ctx context.Context, id bson.ObjectID, title string, content string) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("title", title).Set("content", content).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrap(err, "更新文档失败")
	}
	return nil
}

func (d *DocumentDao) RenameDocumentByID(ctx context.Context, id bson.ObjectID, title string) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Set("title", title).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrap(err, "重命名文档失败")
	}
	return nil
}

func (d *DocumentDao) DeleteByID(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrap(err, "删除文档失败")
	}
	return nil
}

func (d *DocumentDao) DeleteByIDList(ctx context.Context, idList []bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.In("_id", idList...)).DeleteMany(ctx)
	if err != nil {
		return errors.Wrap(err, "删除多个文档失败")
	}
	return nil
}
