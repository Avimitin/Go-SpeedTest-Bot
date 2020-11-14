package speedtest

import (
	"fmt"
	"testing"
)

func TestGetHost(t *testing.T) {
	h := GetHost()
	if h.Port == 0 {
		t.Fail()
	}
}

func TestPing(t *testing.T) {
	fmt.Println(Ping(GetHost()))
}

func TestGetStatus(t *testing.T) {
	fmt.Println(GetStatus(GetHost()))
}

func TestReadSubscriptions(t *testing.T) {
	resp, err := ReadSubscriptions(GetHost(), "")
	if err != nil {
		t.Failed()
	}
	for _, r := range resp {
		fmt.Println(r.Type, r.Config)
	}
}

func TestGetResult(t *testing.T) {
	resp, err := GetResult(GetHost())
	if err != nil {
		t.Failed()
	}
	fmt.Println(resp)
}

func TestStartTest(t *testing.T) {
	subs, err := ReadSubscriptions(GetHost(), "")
	if err != nil {
		t.Failed()
	}
	sCFG := NewStartConfigs("ST_ASYNC", "TCP_PING", subs)
	sCFG.Group = ""
	StartTest(GetHost(), sCFG)
}
