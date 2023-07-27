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
	To(to string, cc ...string)
}

type ByApi interface {
	ByPush
}

type ByPush interface {
	WithBody(body string)
	Send() error
}
