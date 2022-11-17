package twillio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Twillio struct {
	Client *http.Client
	SID    string // twillio's sid
	Token  string // twillio's access token
	URL    string // twillio's post url
	From   string // senders name
}

func New(sid, token, sender string) *Twillio {
	client := &http.Client{}
	t := &Twillio{
		Client: client,
		SID:    sid,
		Token:  token,
		URL:    "https://api.twilio.com/2010-04-01/Accounts/" + sid + "/Messages.json",
		From:   sender,
	}
	return t
}

func sms(ctx context.Context, t *Twillio, To, Body string) error {
	msgData := url.Values{}
	msgData.Set("To", To)
	msgData.Set("From", t.From)
	msgData.Set("Body", Body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", t.URL, &msgDataReader)
	req.SetBasicAuth(t.SID, t.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err := t.Client.Do(req)
	return err
}

func SendPin(ctx context.Context, t *Twillio, To, Pin string) error {
	body := fmt.Sprintf("Your MediaWatch Verification Code is %s", Pin)
	return sms(ctx, t, To, body)
}
