package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	ID         bson.ObjectID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Content    string
	Description  string
	Thumbnail    string
	DocumentType string
	IsPublic     bool
	ParentID     bson.ObjectID
	DocumentID   bson.ObjectID
}
