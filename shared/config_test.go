package shared_test

import (
	"github.com/avelinobego/gestao_pedagogica/shared"
	"testing"
)

func TestConfig(t *testing.T) {
	_, e := shared.LoadConfig("$HOME/.config",
		".")
	if e != nil {
		t.Error(e)
		t.Fail()
	}

}
