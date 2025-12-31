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
	ManageUserEvent    Event = "manage-user"
	ManageNotification Event = "manage-notification"
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
