package dao

import (
	"context"
	"errors"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	mongox.Model `bson:",inline"`
	Title        string
	Content      string
	DocumentType string       `bson:"document_type"`
	IsPublic     bool         `bson:"is_public"`
	ParentID     bson.ObjectID `bson:"parent_id,omitempty"` // 根节点不需要parent_id
	DocumentID   bson.ObjectID `bson:"document_id,omitempty"` // 根节点不需要document_id
}

type IDocumentDao interface {
	Create(ctx context.Context, doc *Document) error
	FindAllRoot(ctx context.Context) ([]*Document, error)
	FindAllByDocumentID(ctx context.Context, documentID bson.ObjectID) ([]*Document, error)
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
			DocumentType: doc.DocumentType,
			IsPublic:     doc.IsPublic,
		})
		if err != nil {
			return err
		}
		if insertResult.InsertedID == nil {
			return errors.New("新增文档失败")
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

func (d *DocumentDao) FindAllRoot(ctx context.Context) ([]*Document, error) {
	documentList, err := d.coll.Finder().Filter(query.NewBuilder().Eq("document_type", "root").Build()).Find(ctx)
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
