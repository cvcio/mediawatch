package link

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func Validate(URL string) (*url.URL, error) {
	l := regexp.MustCompile(`\#..*$`).ReplaceAllString(URL, `$1`)
	l = regexp.MustCompile(`\?utm_source.*$`).ReplaceAllString(l, `$1`)

	if strings.Contains(l, ".jpg") || strings.Contains(l, ".png") {
		return nil, fmt.Errorf("link is an image")
	}
	if strings.Contains(l, ".live24.gr") {
		return nil, fmt.Errorf("link is a radio stream")
	}
	if strings.Contains(l, ".indymedia.org") {
		return nil, fmt.Errorf("link is forbidden")
	}
	// Test URL
	p, err := url.Parse(l)
	// Don't allow Root Domain or Error
	if err != nil {
		return nil, err
	}
	if len(p.Path) <= 1 {
		return nil, fmt.Errorf("link is root domain")
	}

	return p, nil
}

func Parse(URL string) (string, error) {
	p, err := Validate(URL)
	if err != nil {
		return "", err
	}

	var expandedUrl = p.String()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			expandedUrl = req.URL.String()
			return nil
		},
	}

	req, err := http.NewRequest(http.MethodGet, p.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return expandedUrl, nil
}

type CatchedURL struct {
	ID               string `json:"docId"`
	URL              string `json:"url"`
	TweetID          int64  `json:"tweet_id"`
	TwitterUserID    int64  `json:"twitter_user_id"`
	TwitterUserIDStr string `json:"twitter_user_id_str"`
	ScreenName       string `json:"screen_name"`
	CreatedAt        string `json:"created_at"`
	Title            string `json:"title"`
}
