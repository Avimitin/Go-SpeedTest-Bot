package speedtest

import (
	"testing"
)

func NewConfig() []*SubscriptionResp {
	sub := []*SubscriptionResp{
		{
			Config: &ShadowConfig{
				Remarks: "HK",
			},
		},
		{
			Config: &ShadowConfig{
				Remarks: "SG",
			},
		},
		{
			Config: &ShadowConfig{
				Remarks: "US",
			},
		},
	}
	return sub
}

func TestIncludeRemarks(t *testing.T) {
	newcfg := IncludeRemarks(NewConfig(), []string{"HK"})
	if len(newcfg) != 1 {
		t.Fatal("Unexpected array length")
	}
	get := newcfg[0].Config.Remarks
	want := "HK"

	if get != want {
		t.Errorf("get %s want %s", get, want)
	}
}

func TestExcludeRemarks(t *testing.T) {
	newcfg := ExcludeRemarks(NewConfig(), []string{"HK", "US"})
	if len(newcfg) != 1 {
		t.Fatal("Unable array length")
	}
	get := newcfg[0].Config.Remarks
	want := "SG"

	if get != want {
		t.Errorf("get %s want %s", get, want)
	}
}
