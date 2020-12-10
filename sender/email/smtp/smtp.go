package smtp

import (
	"crypto/tls"
	"io"

	gonotifySMTP "gomodules.xyz/notify/smtp"
	"gopkg.in/gomail.v2"
)

const UID = "smtp"

type Options = gonotifySMTP.Options

type client struct {
	opt         Options
	from        string
	subject     string
	body        string
	html        bool
	attachments map[string][]byte
}

func New(opt Options) *client {
	return &client{opt: opt}
}

func (c *client) UID() string {
	return UID
}

func (c *client) From(from string) {
	c.opt.From = from
}

func (c *client) WithSubject(subject string) {
	c.subject = subject
}

func (c *client) WithBody(body string) {
	c.body = body
}

func (c *client) WithTag(tag string) {
}

func (c *client) WithNoTracking() {
}

func (c *client) To(to string, cc ...string) {
	c.opt.To = append([]string{to}, cc...)
}

func (c *client) WithAttachments(attachments map[string][]byte) {
	c.attachments = attachments
}

func (c *client) Send() error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", c.opt.From)
	mail.SetHeader("To", c.opt.To...)
	mail.SetHeader("Subject", c.subject)
	if c.html {
		mail.SetBody("text/html", c.body)
	} else {
		mail.SetBody("text/plain", c.body)
	}

	for filename, filedata := range c.attachments {
		dataCopyFunc := gomail.SetCopyFunc(
			func(data []byte) func(io.Writer) error {
				return func(w io.Writer) error {
					_, err := w.Write(data)
					return err
				}
			}(filedata),
		)
		mail.Attach(filename, dataCopyFunc)
	}

	var d *gomail.Dialer
	if c.opt.Username != "" && c.opt.Password != "" {
		d = gomail.NewDialer(c.opt.Host, c.opt.Port, c.opt.Username, c.opt.Password)
	} else {
		d = &gomail.Dialer{Host: c.opt.Host, Port: c.opt.Port}
	}
	if c.opt.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return d.DialAndSend(mail)
}

func (c *client) SendHtml() error {
	c.html = true
	return c.Send()
}
