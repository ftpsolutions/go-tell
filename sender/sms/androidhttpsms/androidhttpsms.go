package androidhttpsms

import (
	"fmt"
	"net/http"
)

// Built for this application
// https://play.google.com/store/apps/details?id=piersoft.httpsmsserver2

const hTTPSMSURLTemplate = "http://%s:%s/SendSMS/user=%s&password=%s&phoneNumber=%s&msg=%s"

func buildSMSURL(host, port, user, pass, message string, phoneNumbers []string) string {
	numbers := ""
	for _, number := range phoneNumbers {
		numbers += number + ";"
	}
	numbers = numbers[:len(numbers)]

	return fmt.Sprintf(hTTPSMSURLTemplate, host, port, user, pass, numbers, message)
}

type client struct {
	host string
	port string
	user string
	pass string

	from string
	body string
	to   string
	cc   []string
}

func (c *client) UID() string {
	return "SMS"
}

func (c *client) From(from string) {
	c.from = from
}

func (c *client) WithBody(body string) {
	c.body = body
}

func (c *client) To(to string, cc ...string) {
	c.to = to
	c.cc = cc
}

func (c client) Send() error {
	url := buildSMSURL(
		c.host,
		c.port,
		c.user,
		c.pass,
		c.body,
		append([]string{c.to}, c.cc...),
	)

	_, err := http.Get(url)
	if err != nil {
		return err
	}
	return nil
}

func New(host, port, user, pass string) *client {
	return &client{
		host: host,
		port: port,
		user: user,
		pass: pass,
	}
}
