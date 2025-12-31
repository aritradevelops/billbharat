package events

import (
	"time"

	"github.com/google/uuid"
)

type Event string

type EventPayload[T any] struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Event     Event     `json:"event,omitempty"`
	Data      T         `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Action    string    `json:"action,omitempty"`
}

const (
	ManageUserEvent         Event = "manage-user"
	ManageNotification      Event = "manage-notification"
	ManageBusinessEvent     Event = "manage-business"
	ManageBusinessUserEvent Event = "manage-business-user"
)

func newEvent[T any](event Event, action string, data T) EventPayload[T] {
	return EventPayload[T]{
		ID:        uuid.New(),
		Event:     event,
		Data:      data,
		Timestamp: time.Now(),
	}
}

func NewUserManageEvent(action string, data ManageUserEventPayload) EventPayload[ManageUserEventPayload] {
	return newEvent(ManageUserEvent, action, data)
}

func NewBusinessManageEvent(action string, data MangageBusinessEventPayload) EventPayload[MangageBusinessEventPayload] {
	return newEvent(ManageBusinessEvent, action, data)
}

func NewBusinessUserManageEvent(action string, data MangageBusinessUserEventPayload) EventPayload[MangageBusinessUserEventPayload] {
	return newEvent(ManageBusinessUserEvent, action, data)
}

type ManageUserEventPayload struct {
	ID            uuid.UUID  `json:"id"`
	HumanID       string     `json:"human_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Dp            *string    `json:"dp"`
	EmailVerified bool       `json:"email_verified"`
	Phone         string     `json:"phone"`
	PhoneVerified bool       `json:"phone_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by"`
	DeactivatedAt *time.Time `json:"deactivated_at"`
	DeactivatedBy *uuid.UUID `json:"deactivated_by"`
	DeletedAt     *time.Time `json:"deleted_at"`
	DeletedBy     *uuid.UUID `json:"deleted_by"`
}

type MangageBusinessEventPayload struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	Logo            *string    `json:"logo"`
	Industry        string     `json:"industry"`
	PrimaryCurrency string     `json:"primary_currency"`
	OwnerID         uuid.UUID  `json:"owner_id"`
	Currencies      []string   `json:"currencies"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	UpdatedBy       *uuid.UUID `json:"updated_by"`
	DeletedAt       *time.Time `json:"deleted_at"`
	DeletedBy       *uuid.UUID `json:"deleted_by"`
}

type MangageBusinessUserEventPayload struct {
	UserID     uuid.UUID  `json:"user_id"`
	BusinessID uuid.UUID  `json:"business_id"`
	Role       string     `json:"role"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  uuid.UUID  `json:"created_by"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UpdatedBy  *uuid.UUID `json:"updated_by"`
	DeletedAt  *time.Time `json:"deleted_at"`
	DeletedBy  *uuid.UUID `json:"deleted_by"`
}
