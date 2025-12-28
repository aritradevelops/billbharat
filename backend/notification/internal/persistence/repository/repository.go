package repository

import (
	"context"

	"github.com/aritradeveops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradeveops/billbharat/backend/notification/internal/persistence/database"
)

type FindTemplateParams struct {
	Event    dao.EventType   `bson:"event" json:"event"`
	Channel  dao.ChannelType `bson:"channel" json:"channel"`
	Locale   string          `bson:"locale" json:"locale"`
	Scope    string          `bson:"scope" json:"scope"`
	Mimetype string          `bson:"mimetype" json:"mimetype"`
}

type Repository interface {
	CreateTemplate(ctx context.Context, template *dao.Template) (*dao.Template, error)
	FindTemplate(ctx context.Context, params FindTemplateParams) (*dao.Template, error)
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{
		db: db,
	}
}

// CreateTemplate implements Repository.
func (r *repository) CreateTemplate(ctx context.Context, template *dao.Template) (*dao.Template, error) {
	collection := r.db.Collection("templates")
	result, err := collection.InsertOne(ctx, template)
	if err != nil {
		return nil, err
	}
	template.ID = result.InsertedID.(string)
	return template, nil
}

func (r *repository) FindTemplate(ctx context.Context, params FindTemplateParams) (*dao.Template, error) {
	collection := r.db.Collection("templates")

	var template dao.Template
	err := collection.FindOne(ctx, params).Decode(&template)

	if err != nil {
		return nil, err
	}
	return &template, nil
}
