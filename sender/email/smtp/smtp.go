package smtp

import (
	"crypto/tls"
	"io"
	"log"
	"runtime/debug"

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

	for filename, fileData := range c.attachments {
		dataCopyFunc := gomail.SetCopyFunc(
			func(data []byte) func(io.Writer) error {
				return func(w io.Writer) error {
					_, err := w.Write(data)
					return err
				}
			}(fileData),
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

	// TODO: the panic handling below is to deal with https://github.com/go-gomail/gomail/issues/89

	var err error

	// this will get called whenever we're leaving the scope of this function
	defer func() {

		// check if we're leaving the scope of the function because of a panic
		r := recover()

		// if we are, log some information to the user
		if r != nil {

			// set the error to be returned from the outermost function
			err = r.(error)

			// logs for days
			log.Printf("error: caught and handled panic; err=%#+v; see stack trace below", err)
			log.Printf("")
			log.Printf(">>>> caught stack trace")
			log.Printf(string(debug.Stack()))
			log.Printf("<<<< caught stack trace")
			log.Printf("")
		}
	}()

	// this is the call that may panic
	err = d.DialAndSend(mail)

	// in which case, this err will be the one learned of in the recover
	return err
}

func (c *client) SendHtml() error {
	c.html = true
	return c.Send()
}
