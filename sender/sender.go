package sender

type ByEmail interface {
	BySMS
	WithSubject(subject string)
	WithTag(tag string)
	WithNoTracking()
	SendHtml() error
}

type ByEmailWithAttachments interface {
	ByEmail
	WithAttachments(attachments map[string][]byte)
}

type BySMS interface {
	From(from string)
	ByChat
}

type ByChat interface {
	ByPush
}

type ByPush interface {
	UID() string
	WithBody(body string)
	To(to string, cc ...string)
	Send() error
}
