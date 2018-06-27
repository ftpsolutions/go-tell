package sender

import notify "github.com/appscode/go-notify"

type goNotifyEmailWrapper struct {
	email notify.ByEmail
}
type goNotifySMSWrapper struct {
	sms notify.BySMS
}
type goNotifyChatWrapper struct {
	chat notify.ByChat
}
type goNotifyPushWrapper struct {
	push notify.ByPush
}

func (s *goNotifyEmailWrapper) UID() string { return s.email.UID() }
func (s *goNotifySMSWrapper) UID() string   { return s.sms.UID() }
func (s *goNotifyChatWrapper) UID() string  { return s.chat.UID() }
func (s *goNotifyPushWrapper) UID() string  { return s.push.UID() }

func (s *goNotifyEmailWrapper) WithBody(body string) { s.email = s.email.WithBody(body) }
func (s *goNotifySMSWrapper) WithBody(body string)   { s.sms = s.sms.WithBody(body) }
func (s *goNotifyChatWrapper) WithBody(body string)  { s.chat = s.chat.WithBody(body) }
func (s *goNotifyPushWrapper) WithBody(body string)  { s.push = s.push.WithBody(body) }

func (s *goNotifyEmailWrapper) To(to string, cc ...string) { s.email = s.email.To(to, cc...) }
func (s *goNotifySMSWrapper) To(to string, cc ...string)   { s.sms = s.sms.To(to, cc...) }
func (s *goNotifyChatWrapper) To(to string, cc ...string)  { s.chat = s.chat.To(to, cc...) }
func (s *goNotifyPushWrapper) To(to string, cc ...string) {
	s.push = s.push.To(append([]string{to}, cc...)...)
}

func (s *goNotifyEmailWrapper) Send() error { return s.email.Send() }
func (s *goNotifySMSWrapper) Send() error   { return s.sms.Send() }
func (s *goNotifyChatWrapper) Send() error  { return s.chat.Send() }
func (s *goNotifyPushWrapper) Send() error  { return s.push.Send() }

func (s *goNotifyEmailWrapper) From(from string) { s.email = s.email.From(from) }
func (s *goNotifySMSWrapper) From(from string)   { s.sms = s.sms.From(from) }

func (s *goNotifyEmailWrapper) WithSubject(subject string) { s.email = s.email.WithSubject(subject) }
func (s *goNotifyEmailWrapper) WithTag(tag string)         { s.email = s.email.WithTag(tag) }
func (s *goNotifyEmailWrapper) WithNoTracking()            { s.email = s.email.WithNoTracking() }
func (s *goNotifyEmailWrapper) SendHtml() error            { return s.email.SendHtml() }

func WrapGoNotifyEmail(email notify.ByEmail) ByEmail {
	return &goNotifyEmailWrapper{email}
}

func WrapGoNotifySMS(sms notify.BySMS) BySMS {
	return &goNotifySMSWrapper{sms}
}

func WrapGoNotifyChat(chat notify.ByChat) ByChat {
	return &goNotifyChatWrapper{chat}
}

func WrapGoNotifyPush(push notify.ByPush) ByPush {
	return &goNotifyPushWrapper{push}
}
