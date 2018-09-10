package sender

import "fmt"

type Stdout struct {
	uID         string
	from        string
	subject     string
	to          []string
	body        string
	tag         string
	tracking    bool
	html        bool
	attachments map[string][]byte
}

func (s *Stdout) UID() string {
	return s.uID
}

func (s *Stdout) From(from string) {
	s.from = from
}

func (s *Stdout) WithSubject(subject string) {
	s.subject = subject
}

func (s *Stdout) WithBody(body string) {
	s.body = body
}

func (s *Stdout) WithAttachments(attachments map[string][]byte) {
	s.attachments = attachments
}

func (s *Stdout) WithTag(tag string) {
	s.tag = tag
}

func (s *Stdout) WithNoTracking() {
	s.tracking = false
}

func (s *Stdout) To(to string, cc ...string) {
	s.to = []string{to}
	s.to = append(s.to, cc...)
}

func (s *Stdout) Send() error {
	fmt.Println(s)
	return nil
}

func (s *Stdout) SendHtml() error {
	s.html = true
	fmt.Println(s)
	return nil
}

func NewStdout() *Stdout {
	return &Stdout{
		tracking: true,
	}
}
