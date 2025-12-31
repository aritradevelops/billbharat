package events

import (
	"context"
)

type EventManager interface {
	EmitManageUserEvent(ctx context.Context, data EventPayload[ManageUserEventPayload]) error
	OnManageUserEvent(ctx context.Context, handler func(EventPayload[ManageUserEventPayload]) error)
	EmitManageNotificationEvent(ctx context.Context, data EventPayload[MangageNotificationEventPayload]) error
	OnManageNotificationEvent(ctx context.Context, handler func(EventPayload[MangageNotificationEventPayload]) error)
	EmitManageBusinessEvent(ctx context.Context, data EventPayload[MangageBusinessEventPayload]) error
	OnManageBusinessEvent(ctx context.Context, handler func(EventPayload[MangageBusinessEventPayload]) error)
	EmitManageBusinessUserEvent(ctx context.Context, data EventPayload[MangageBusinessUserEventPayload]) error
	OnManageBusinessUserEvent(ctx context.Context, handler func(EventPayload[MangageBusinessUserEventPayload]) error)
}
