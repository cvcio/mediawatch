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

type El_Efsyn struct{}

func (h El_Efsyn) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	url := "https://www.efsyn.gr/news"

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

	collection := doc.Find("article.default-teaser__article")
	if collection.Size() == 0 {
		return list, nil
	}

	collection.Each(func(i int, s *goquery.Selection) {
		href, hok := s.Find("a.full-link").Attr("href")
		if !hok {
			return
		}
		datetime := s.Find("time.default-teaser__date").Text()
		t, err := time.ParseInLocation("02.01.2006, 15:04", strings.TrimSpace(strings.ReplaceAll(datetime, "\n", "")), loc)
		if err != nil {
			return
		}
		item := &gofeed.Item{
			Title:           strings.TrimSpace(s.Find("div.default-teaser__title").Text()),
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            "https://www.efsyn.gr" + href,
		}
		list = append(list, item)
	})

	return list, nil
}
