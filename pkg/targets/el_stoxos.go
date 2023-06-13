package targets

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/cvcio/mediawatch/pkg/helper"
	"github.com/mmcdole/gofeed"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type El_Stoxos struct{}

func (h El_Stoxos) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	dates := map[string]string{
		"ΙΑΝΟΥΑΡΙΟΥ":  "01",
		"ΦΕΒΡΟΥΑΡΙΟΥ": "02",
		"ΜΑΡΤΙΟΥ":     "03",
		"ΑΠΡΙΛΙΟΥ":    "04",
		"ΜΑΙΟΥ":       "05",
		"ΙΟΥΝΙΟΥ":     "06",
		"ΙΟΥΛΙΟΥ":     "07",
		"ΑΥΓΟΥΣΤΟΥ":   "08",
		"ΣΕΠΤΕΜΒΡΙΟΥ": "09",
		"ΟΚΤΩΒΡΙΟΥ":   "10",
		"ΝΟΕΜΒΡΙΟΥ":   "11",
		"ΔΕΚΕΜΒΡΙΟΥ":  "12",
	}

	url := "https://www.stoxos.gr/search?max-results=24"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Referer", "https://news.google.com/")
	req.Header.Set("User-Agent", helper.RandomUserAgent())

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Failed with code: %d", r.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse document body: %s", err))
	}

	var list []*gofeed.Item
	loc, _ := time.LoadLocation("Europe/Athens")
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	collection := doc.Find("article")
	if collection.Size() == 0 {
		return list, nil
	}

	collection.Each(func(i int, s *goquery.Selection) {
		href, hok := s.Find("h3.post-title a").Attr("href")
		if !hok {
			return
		}
		datetime := s.Find(".post-timestamp").Text()
		tr := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
		datetime, _, _ = transform.String(tr, datetime)
		datetime = strings.TrimSpace(datetime)
		datetime = strings.ToUpper(datetime)
		for k, v := range dates {
			datetime = strings.ReplaceAll(datetime, k, v)
		}
		t, err := time.ParseInLocation("01 02, 2006", datetime, loc)
		if err != nil {
			return
		}
		now := time.Now().In(loc)
		published := time.Date(t.Year(), t.Month(), t.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), loc)
		item := &gofeed.Item{
			Title:           strings.TrimSpace(s.Find("h3.post-title").Text()),
			Published:       published.Format(time.RFC3339),
			PublishedParsed: &published,
			Link:            href,
		}
		list = append(list, item)
	})

	return list, nil
}
