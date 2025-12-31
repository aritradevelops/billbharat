package notification

type EmailData struct {
	To  []string `json:"to"`
	CC  []string `json:"cc"`
	BCC []string `json:"bcc"`
}

func NewEmail(to ...string) *EmailData {
	return &EmailData{
		To:  to,
		CC:  []string{},
		BCC: []string{},
	}
}

func (e *EmailData) WithCC(cc ...string) *EmailData {
	e.CC = append(e.CC, cc...)
	return e
}

func (e *EmailData) WithBCC(bcc ...string) *EmailData {
	e.BCC = append(e.BCC, bcc...)
	return e
}
