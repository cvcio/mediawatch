package targets_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/proxy"
	"github.com/cvcio/mediawatch/pkg/targets"
	"github.com/kelseyhightower/envconfig"
)

var client_no_proxy = &http.Client{Timeout: 30 * time.Second}

// Amna
func Test_El_Amna_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "amna.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Amna(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "amna.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// News247
func Test_El_News247_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "news247.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))
	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_News247(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "news247.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))
	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// ProNews
func Test_El_ProNews_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "pronews.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_ProNews(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "pronews.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// Liberal
func Test_El_Liberal_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "liberal.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Liberal(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "liberal.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// Lifo
func Test_El_Lifo_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "lifo.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Lifo(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "lifo.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}

}

// Efsyn
func Test_El_Efsyn_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "efsyn.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Efsyn(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "efsyn.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// MoneyReview
func Test_El_MoneyReview_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "moneyreview.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_MoneyReview(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "moneyreview.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// Stoxos
func Test_El_Stoxos_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "stoxos.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Stoxos(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "stoxos.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

// Tvxs
func Test_El_Tvxs_proxy_off(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "tvxs.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(client_no_proxy, client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}

func Test_El_Tvxs(t *testing.T) {
	cfg := config.NewConfig()
	_ = envconfig.Process("", cfg)

	hostname := "tvxs.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyURL()), client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't Get document for %s, err: %s", hostname, err)
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Logf("%s %s %s", v.Published, v.Title, v.Link)
	}
}
