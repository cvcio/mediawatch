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

type El_Tvxs struct{}

func (h El_Tvxs) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	url := "https://tvxs.gr/latest-news/"

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

	collection := doc.Find("article.col2of3:not(:first-child)")
	if collection.Size() == 0 {
		return list, nil
	}

	collection.Each(func(i int, s *goquery.Selection) {
		href, hok := s.Find("div.col1of3.trim-left-a a").Attr("href")
		if !hok {
			return
		}
		d1, dok := s.Find("time").Attr("datetime")
		if !dok {
			return
		}
		t1, err := time.ParseInLocation("2006-01-02 15:04", strings.TrimSpace(d1), loc)
		if err != nil {
			return
		}

		d2 := s.Find("time").Text()
		t2, err := time.ParseInLocation("15:04", strings.TrimSpace(d2), loc)
		if err != nil {
			return
		}

		t := time.Date(t1.Year(), t1.Month(), t1.Day(), t2.Hour(), t2.Minute(), 0, 0, loc)
		item := &gofeed.Item{
			Title:           strings.TrimSpace(s.Find("div.col1of3 h3").Text()),
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            "https:" + href,
		}
		list = append(list, item)
	})

	return list, nil
}
