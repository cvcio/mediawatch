package mailer

import (
	"bytes"
	"crypto/tls"
	"html/template"

	"gopkg.in/gomail.v2"
)

// Mailer struct.
type Mailer struct {
	Dialer   *gomail.Dialer
	BaseURL  string
	From     string // no-reply@mediawatch.io
	FromName string // MediaWatch
	ReplyTo  string // press@mediawatch.io
}

// NewMailer returns a new gomail client.
func NewMailer(smtp string, port int, username, password string, baseURL, from, fromName, replyTo string) *Mailer {
	dialer := gomail.NewDialer(smtp, port, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Mailer{
		dialer,
		baseURL,
		from,
		fromName,
		replyTo,
	}
}

func (m *Mailer) message(to, subject, body string) error {
	gm := gomail.NewMessage()
	gm.SetAddressHeader("From", m.From, m.FromName)
	gm.SetHeader("To", to)
	gm.SetHeader("Reply-To", m.ReplyTo)
	gm.SetHeader("Subject", subject)
	gm.SetBody("text/html", body)

	return m.Dialer.DialAndSend(gm)
}

// SendInvitation to a user.
func (m *Mailer) SendInvitation() error {
	return nil
}

// SendMagicLink to a user.
func (m *Mailer) SendMagicLink(to, subject string, options map[string]interface{}) error {
	options["baseURL"] = m.BaseURL

	var tpl bytes.Buffer
	tmpl, err := template.New("SendMagicLink").Parse(TemplateMagicLink)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(&tpl, options); err != nil {
		return err
	}

	return m.message(to, subject, tpl.String())
}

// SendOTP to a user.
func (m *Mailer) SendOTP() error {
	return nil
}

// SendMessage to a user.
func (m *Mailer) SendMessage() error {
	return nil
}

// package mailer

// import (
// 	"bytes"
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"strings"
// 	"text/template"

// 	gomail "gopkg.in/gomail.v2"
// )

// type Mailer struct {
// 	Dialer   *gomail.Dialer
// 	From     string // no-reply@mediawatch.io
// 	FromName string // MediaWatch
// 	ReplyTo  string // press@mediawatch.io
// }

// func New(smtp string, port int, username, password string, From, FromName, ReplyTo string) *Mailer {
// 	d := gomail.NewDialer(smtp, port, username, password)
// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	m := &Mailer{Dialer: d}
// 	m.From = From
// 	m.FromName = FromName
// 	m.ReplyTo = ReplyTo
// 	return m
// }

// func message(ctx context.Context, m *Mailer, To, subject, body string) error {
// 	gm := gomail.NewMessage()
// 	gm.SetAddressHeader("From", m.From, m.FromName)
// 	gm.SetHeader("To", To)
// 	gm.SetHeader("Reply-To", m.ReplyTo)
// 	gm.SetHeader("Subject", subject)
// 	gm.SetBody("text/html", body)
// 	return m.Dialer.DialAndSend(gm)
// }

// func MessageSimple(ctx context.Context, m *Mailer, To, subject, body string) error {
// 	msgBody := fmt.Sprintf(msgDefault, body)
// 	return message(ctx, m, To, subject, msgBody)
// }

// func SendInvite(ctx context.Context, m *Mailer, To, First, Last, email, nonce, team string) error {
// 	options := map[string]interface{}{}

// 	options["name"] = fmt.Sprintf("%s %s", First, Last)
// 	options["email"] = strings.ToLower(email)
// 	options["nonce"] = nonce
// 	options["team"] = team

// 	var tpl bytes.Buffer
// 	tmpl, err := template.New("SendInvite").Parse(msgInvitation)
// 	if err != nil {
// 		return err
// 	}
// 	if err := tmpl.Execute(&tpl, options); err != nil {
// 		return err
// 	}

// 	return message(ctx, m, To, fmt.Sprintf("Join MediaWatch (Invitation by %s)", options["name"]), tpl.String())
// }

// func SendInviteExistingUser(ctx context.Context, m *Mailer, To, First, Last, email, nonce, team, orgId, memberId string) error {
// 	options := map[string]interface{}{}

// 	options["name"] = fmt.Sprintf("%s %s", First, Last)
// 	options["email"] = strings.ToLower(email)
// 	options["nonce"] = nonce
// 	options["team"] = team
// 	options["orgId"] = orgId
// 	options["memberId"] = memberId

// 	var tpl bytes.Buffer
// 	tmpl, err := template.New("SendInvite").Parse(msgInvitationExistingAccount)
// 	if err != nil {
// 		return err
// 	}
// 	if err := tmpl.Execute(&tpl, options); err != nil {
// 		return err
// 	}
// 	return message(ctx, m, To, fmt.Sprintf("Join MediaWatch (Invitation by %s)", options["name"]), tpl.String())
// }

// func SendNewPass(ctx context.Context, m *Mailer, To, First, pass string) error {
// 	msgBody := fmt.Sprintf(msgNewPass, First, pass)
// 	return message(ctx, m, To, "Password Reset", msgBody)
// }

// func SendReset(ctx context.Context, m *Mailer, To, First, pin, id string) error {
// 	msgBody := fmt.Sprintf(msgReset, First, pin, id, id)
// 	return message(ctx, m, To, "Your Verification Code", msgBody)
// }

// func SendPin(ctx context.Context, m *Mailer, To, First, pin, id string) error {
// 	msgBody := fmt.Sprintf(msgPin, First, pin, id, id)
// 	return message(ctx, m, To, "Your Verification Code", msgBody)
// }

// func SendAccountDeletion(ctx context.Context, m *Mailer, To, First string) error {
// 	msgBody := fmt.Sprintf(msgAccountDeletion, First)
// 	return message(ctx, m, To, "MediaWatch Account Removal", msgBody)
// }
