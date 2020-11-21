package speedtest

import (
	"fmt"
	"os"
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
	subs, err := ReadSubscriptions(GetHost(), "https://oxygenproxy.com/auth/register")
	if err != nil {
		t.Failed()
	}
	sCFG := NewStartConfigs("ST_ASYNC", "TCP_PING", subs)
	StartTest(GetHost(), sCFG, make(chan string))
}

func NewConfig() []*SubscriptionResp {
	sub, err := ReadSubscriptions(GetHost(), "")
	if err != nil {
		os.Exit(-1)
	}
	return sub
}

func TestIncludeRemarks(t *testing.T) {
	newcfg := IncludeRemarks(NewConfig(), []string{"香港"})
	fmt.Println(newcfg)
}

func TestExcludeRemarks(t *testing.T) {
	newcfg := ExcludeRemarks(NewConfig(), []string{"剩余", "台湾", "香港", "过期"})
	for _, n := range newcfg {
		fmt.Printf("%+v\n", n.Config.Remarks)
	}
}
