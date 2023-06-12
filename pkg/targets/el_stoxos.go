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

type El_Stoxos struct{}

func (h El_Stoxos) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	loc, _ := time.LoadLocation("Europe/Athens")

	now := time.Now().In(loc)
	url := "https://www.liberal.gr/news-feed/" + now.Format("20060102")

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
	collection := doc.Find(".listing #block-liberal-theme-content .article")
	if collection.Size() == 0 {
		return list, nil
	}

	collection.Each(func(i int, s *goquery.Selection) {
		href, hok := s.Find("a").Attr("href")
		if !hok {
			return
		}
		datetime := s.Find(".article__info").Text()
		datetime = strings.TrimSpace(datetime)
		datetime = strings.Split(datetime, "\n")[0]
		t, err := time.ParseInLocation("02/01/2006 • 15:04", strings.TrimSpace(datetime), loc)
		if err != nil {
			return
		}
		item := &gofeed.Item{
			Title:           strings.TrimSpace(s.Find(".article__title p").Text()),
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            "https://www.liberal.gr" + href,
		}
		list = append(list, item)
	})

	return list, nil
}
