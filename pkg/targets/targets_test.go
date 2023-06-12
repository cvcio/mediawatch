package targets_test

import (
	"testing"

	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/proxy"
	"github.com/cvcio/mediawatch/pkg/targets"
	"github.com/kelseyhightower/envconfig"
)

func Test_El_Amna(t *testing.T) {
	cfg := config.NewConfig()
	envconfig.Process("", cfg)

	hostname := "amna.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))

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
	envconfig.Process("", cfg)

	hostname := "news247.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))
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
	envconfig.Process("", cfg)

	hostname := "pronews.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))

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
	envconfig.Process("", cfg)

	hostname := "liberal.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))

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
	envconfig.Process("", cfg)

	hostname := "lifo.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))

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
	envconfig.Process("", cfg)

	hostname := "efsyn.gr"
	client := targets.Targets[hostname]

	res, err := targets.ParseList(proxy.CreateProxy(cfg.GetProxyList(), cfg.Proxy.UserName, cfg.Proxy.Password), client.(targets.Target))

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
