package speedtest

import (
	"testing"
)

func TestGetHost(t *testing.T) {
	h := GetHost()
	if h.Port == 0 {
		t.Fail()
	}
}
