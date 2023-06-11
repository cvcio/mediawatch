package targets

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

// Ex. https://www.amna.gr/sport/article/737728/I-kardia-tou-Mpasket-tha-chtupa-sto-Irakleio-Kritisrn
type AmnaGRItem struct {
	Id        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Category  string `json:"parent_title,omitempty"`
	Note      string `json:"note2,omitempty"`
	Path      string `json:"note3,omitempty"`
	Published string `json:"c_daytime,omitempty"` // 2023-06-11 10:20:36
	Kind      string `json:"kind,omitempty"`
}

type AmnaGR struct{}

func (h AmnaGR) Get() ([]*gofeed.Item, error) {
	r, err := http.Get("https://www.amna.gr/feeds/getfolder.php?id=46&infolevel=INTERMEDIATE&offset=0&numrows=30&kind=article&byrole=false&subfolders=true&order=[[%22c_timestamp%22,%22desc%22]]&exclude=")
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	var res []*AmnaGRItem
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	var list []*gofeed.Item

	for _, item := range res {
		t, _ := time.Parse(time.DateTime, item.Published)
		list = append(list, &gofeed.Item{
			Title:           item.Title,
			Published:       t.Format(time.RFC3339),
			PublishedParsed: &t,
			Link:            "https://www.amna.gr/feeds/getarticle.php?id=" + item.Id,
			// Link:            "https://www.amna.gr/" + item.Note + "/" + item.Kind + "/" + item.Id + "/" + item.Path,
		})
	}

	return list, nil
}
