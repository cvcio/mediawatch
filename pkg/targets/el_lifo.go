package targets

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cvcio/mediawatch/pkg/helper"
	"github.com/mmcdole/gofeed"
)

type El_Lifo struct{}

func (h El_Lifo) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	url := "https://www.lifo.gr/now"

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
	collection := doc.Find(".view-rows article")
	if collection.Size() == 0 {
		return list, nil
	}

	collection.Each(func(i int, s *goquery.Selection) {
		href, hok := s.Find("h3 a:last-of-type").Attr("href")
		if !hok {
			return
		}
		datetime, dok := s.Find("time").Attr("datetime")
		if !dok {
			return
		}

		t, err := time.Parse("2006-01-02T15:04:05-0700", datetime)
		if err != nil {
			return
		}
		title := s.Find("h3 a:last-of-type").Text()
		item := &gofeed.Item{
			Title:           strings.TrimSpace(title),
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            href,
		}
		list = append(list, item)
	})

	return list, nil
}
