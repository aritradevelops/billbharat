package notification

type Event string

const (
	SIGNUP             Event = "signup"
	EMAIL_VERIFICATION Event = "email_verification"
	EMAIL_VERIFIED     Event = "email_verified"
	PHONE_VERIFICATION Event = "phone_verification"
	PHONE_VERIFIED     Event = "phone_verified"
	FORGOT_PASSWORD    Event = "forgot_password"
	RESET_PASSWORD     Event = "reset_password"
	CHANGE_PASSWORD    Event = "change_password"
	USER_INVITED       Event = "user_invited"
)

type Channel string

const (
	EMAIL    Channel = "email"
	SMS      Channel = "sms"
	PUSH     Channel = "push"
	WHATSAPP Channel = "whatsapp"
)

type Kind string

const (
	P2P       Kind = "p2p"
	BROADCAST Kind = "broadcast"
)

func NewSMS(to ...string) *SMSData {
	return &SMSData{
		To: to,
	}
}

type SMSData struct {
	To []string `json:"to"`
}

func NewPushMessage(to ...string) *PushMessageData {
	return &PushMessageData{
		To: to,
	}
}

type PushMessageData struct {
	To []string `json:"to"`
}

type WhatsappMessageData struct {
	To []string `json:"to"`
}

func NewWhatsappMessage(to ...string) *WhatsappMessageData {
	return &WhatsappMessageData{
		To: to,
	}
}
