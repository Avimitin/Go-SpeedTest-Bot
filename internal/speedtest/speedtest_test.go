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
