package repository

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FindTemplateParams struct {
	Event    notification.Event   `bson:"event" json:"event"`
	Channel  notification.Channel `bson:"channel" json:"channel"`
	Locale   string               `bson:"locale" json:"locale"`
	Scope    string               `bson:"scope" json:"scope"`
	Mimetype string               `bson:"mimetype" json:"mimetype"`
}

type Repository interface {
	CreateTemplate(ctx context.Context, template dao.Template) (dao.Template, error)
	FindTemplate(ctx context.Context, params FindTemplateParams) (dao.Template, error)
	SyncUser(ctx context.Context, user dao.User) error
	SyncBusiness(ctx context.Context, business dao.Business) error
	SyncBusinessUser(ctx context.Context, businessUser dao.BusinessUser) error
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
func (r *repository) CreateTemplate(ctx context.Context, template dao.Template) (dao.Template, error) {
	collection := r.db.Collection("templates")
	result, err := collection.InsertOne(ctx, template)
	if err != nil {
		return dao.Template{}, err
	}
	template.ID = result.InsertedID.(uuid.UUID)
	return template, nil
}

func (r *repository) FindTemplate(ctx context.Context, params FindTemplateParams) (dao.Template, error) {
	collection := r.db.Collection("templates")

	var template dao.Template
	err := collection.FindOne(ctx, params).Decode(&template)

	if err != nil {
		return dao.Template{}, err
	}
	return template, nil
}

func (r *repository) SyncUser(ctx context.Context, user dao.User) error {
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SyncBusiness(ctx context.Context, business dao.Business) error {
	collection := r.db.Collection("businesses")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": business.ID}, bson.M{"$set": business}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SyncBusinessUser(ctx context.Context, businessUser dao.BusinessUser) error {
	collection := r.db.Collection("business_users")
	_, err := collection.UpdateOne(ctx, bson.M{"user_id": businessUser.UserID, "business_id": businessUser.BusinessID}, bson.M{"$set": businessUser}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
