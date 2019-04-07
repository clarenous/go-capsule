package api

import (
	_ "github.com/clarenous/go-capsule/consensus/algorithm/pow"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	a := NewAPI()

	err := a.Start()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(time.Second * 600)

	a.Stop()
}
