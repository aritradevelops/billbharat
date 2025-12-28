package dao

import "github.com/google/uuid"

type EventType string

const (
	USER_SIGNUP_OTP EventType = "USER_SIGNUP_OTP"
)

type ChannelType string

const (
	EMAIL ChannelType = "EMAIL"
	SMS   ChannelType = "SMS"
	PUSH  ChannelType = "PUSH"
)

type Template struct {
	ID string `bson:"_id" json:"_id"`
	// if later has to move to postgres
	UID      uuid.UUID   `bson:"id" json:"id"`
	Event    EventType   `bson:"event" json:"event"`
	Channel  ChannelType `bson:"channel" json:"channel"`
	Locale   string      `bson:"locale" json:"locale"`
	Subject  string      `bson:"subject" json:"subject"`
	Body     string      `bson:"body" json:"body"`
	Mimetype string      `bson:"mimetype" json:"mimetype"`
	// scope can be default or some custom target like for a business the business id
	Scope string `bson:"scope" json:"scope"`
	// we don't need version as of now
	// Version string `bson:"version" json:"version"`
}
