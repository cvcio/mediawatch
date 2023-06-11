package targets

import "github.com/mmcdole/gofeed"

type Target interface {
	Get() ([]*gofeed.Item, error)
}

func Get(t Target) ([]*gofeed.Item, error) {
	return t.Get()
}

var Targets = map[string]interface{}{
	"amna.gr": AmnaGR{},
}
