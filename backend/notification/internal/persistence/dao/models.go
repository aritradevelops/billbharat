package dao

import (
	"time"

	"github.com/aritradevelops/billbharat/backend/shared/notification"
	"github.com/google/uuid"
)

type Template struct {
	ID       uuid.UUID                      `bson:"_id" json:"id"`
	Event    notification.NotificationEvent `bson:"event" json:"event"`
	Channel  notification.ChannelType       `bson:"channel" json:"channel"`
	Locale   string                         `bson:"locale" json:"locale"`
	Subject  string                         `bson:"subject" json:"subject"`
	Body     string                         `bson:"body" json:"body"`
	Mimetype string                         `bson:"mimetype" json:"mimetype"`
	// scope can be default or some custom target like for a business the business id
	Scope string `bson:"scope" json:"scope"`
	// we don't need version as of now
	// Version string `bson:"version" json:"version"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	CreatedBy uuid.UUID  `bson:"created_by" json:"created_by"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	UpdatedBy *uuid.UUID `bson:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at"`
	DeletedBy *uuid.UUID `bson:"deleted_by" json:"deleted_by"`
}

type User struct {
	ID            uuid.UUID  `bson:"_id" json:"id"`
	HumanID       string     `bson:"human_id" json:"human_id"`
	Name          string     `bson:"name" json:"name"`
	Email         string     `bson:"email" json:"email"`
	Dp            *string    `bson:"dp" json:"dp"`
	EmailVerified bool       `bson:"email_verified" json:"email_verified"`
	Phone         string     `bson:"phone" json:"phone"`
	PhoneVerified bool       `bson:"phone_verified" json:"phone_verified"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	CreatedBy     uuid.UUID  `bson:"created_by" json:"created_by"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
	UpdatedBy     *uuid.UUID `bson:"updated_by" json:"updated_by"`
	DeactivatedAt *time.Time `bson:"deactivated_at" json:"deactivated_at"`
	DeactivatedBy *uuid.UUID `bson:"deactivated_by" json:"deactivated_by"`
	DeletedAt     *time.Time `bson:"deleted_at" json:"deleted_at"`
	DeletedBy     *uuid.UUID `bson:"deleted_by" json:"deleted_by"`
}
