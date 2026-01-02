package events

import "github.com/aritradevelops/billbharat/backend/shared/notification"

type ManageNotificationEventPayload struct {
	Event   notification.Event           `json:"event"`
	Kind    notification.Kind            `json:"kind"`
	Payload []NotificationChannelPayload `json:"payload"`
	Tokens  any                          `json:"tokens"`
}

type NotificationChannelPayload struct {
	Channel notification.Channel `json:"channel"`
	Data    any                  `json:"data"`
}

func NewNotificationManageEvent(data ManageNotificationEventPayload) EventPayload[ManageNotificationEventPayload] {
	return newEvent(ManageNotification, "send", data)
}
