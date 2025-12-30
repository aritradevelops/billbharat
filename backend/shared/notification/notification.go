package notification

type EventType string

const (
	P2P       EventType = "p2p"
	Broadcast EventType = "broadcast"
)

type ChannelType string

const (
	EMAIL ChannelType = "email"
	SMS   ChannelType = "sms"
	PUSH  ChannelType = "push"
)

type Channel[Data any] struct {
	Type ChannelType `json:"type,omitempty"`
	Data Data        `json:"data,omitempty"`
}

type NotificationEvent string

const (
	SignupNotificationEvent NotificationEvent = "signup"
	EmailVerificationEvent  NotificationEvent = "email_verification"
	EmailVerifiedEvent      NotificationEvent = "email_verified"
	PhoneVerificationEvent  NotificationEvent = "phone_verification"
	PhoneVerifiedEvent      NotificationEvent = "phone_verified"
	ForgotPasswordEvent     NotificationEvent = "forgot_password"
	ResetPasswordEvent      NotificationEvent = "reset_password"
	ChangePasswordEvent     NotificationEvent = "change_password"
)

type Notification struct {
	Event    NotificationEvent `json:"event"`
	Channels []Channel[any]    `json:"channels"`
	Data     any               `json:"data"`
}
