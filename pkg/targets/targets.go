package targets

import (
	"net/http"

	"github.com/mmcdole/gofeed"
)

type Target interface {
	ParseList(client *http.Client) ([]*gofeed.Item, error)
}

func ParseList(client *http.Client, t Target) ([]*gofeed.Item, error) {
	return t.ParseList(client)
}

var Targets = map[string]interface{}{
	"amna.gr":        El_Amna{},
	"efsyn.gr":       El_Efsyn{},
	"liberal.gr":     El_Liberal{},
	"lifo.gr":        El_Lifo{},
	"news247.gr":     El_News247{},
	"pronews.gr":     El_ProNews{},
	"moneyreview.gr": El_MoneyReview{},
	"stoxos.gr":      El_Stoxos{},
	"tvxs.gr":        El_Tvxs{},
}
