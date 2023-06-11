package targets_test

import (
	"testing"

	"github.com/cvcio/mediawatch/pkg/targets"
)

func Test_AmnaGR(t *testing.T) {
	hostname := "amna.gr"

	// client := targets.AmnaGR{}

	client := targets.Targets[hostname]
	res, err := targets.Get(client.(targets.Target))

	if err != nil {
		t.Fatalf("Couldn't create Twitter API HTTP Client")
	}

	if len(res) == 0 {
		t.Fail()
	}

	for _, v := range res {
		t.Log(v)
	}
}
