package targets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// Ex. https://www.amna.gr/sport/article/737728/I-kardia-tou-Mpasket-tha-chtupa-sto-Irakleio-Kritisrn
type El_Amna_Item struct {
	Id        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Category  string `json:"parent_title,omitempty"`
	Note      string `json:"note2,omitempty"`
	Path      string `json:"note3,omitempty"`
	Published string `json:"c_daytime,omitempty"` // 2023-06-11 10:20:36
	Kind      string `json:"kind,omitempty"`
}

type El_Amna struct{}

func (h El_Amna) ParseList(client *http.Client) ([]*gofeed.Item, error) {
	r, err := http.Get("https://www.amna.gr/feeds/getfolder.php?id=46&infolevel=INTERMEDIATE&offset=0&numrows=30&kind=article&byrole=false&subfolders=true&order=[[%22c_timestamp%22,%22desc%22]]&exclude=")
	if err != nil {
		return nil, err
	}

	defer func() { _ = r.Body.Close() }()

	if r.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Failed with code: %d", r.StatusCode))
	}

	var res []*El_Amna_Item
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	var list []*gofeed.Item
	loc, _ := time.LoadLocation("Europe/Athens")

	for _, item := range res {
		t, _ := time.ParseInLocation(time.DateTime, item.Published, loc)
		list = append(list, &gofeed.Item{
			Title:           strings.TrimSpace(item.Title),
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            "https://www.amna.gr/feeds/getarticle.php?id=" + item.Id,
		})
	}

	return list, nil
}
