package speedtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-speedtest-bot/module/runner"
	"go-speedtest-bot/module/web"
	"strings"
)

// Ping will test if the given host is accessible or not
func Ping(r runner.Runner) bool {
	resp, err := web.Get(r.Host.GetURL() + "/" + "getversion")
	if err != nil {
		return false
	}
	var v Version
	err = json.Unmarshal(resp, &v)
	if err != nil {
		return false
	}
	return v.Main != "" && v.WebAPI != ""
}

// GetStatus is used for fetching backend status
func GetStatus(r runner.Runner) (*Status, error) {
	resp, err := web.Get(r.Host.GetURL() + "/" + "status")
	if err != nil {
		return nil, fmt.Errorf("get status: %v", err)
	}
	var st *Status
	err = json.Unmarshal(resp, &st)
	if err != nil {
		return nil, fmt.Errorf("decode %q: %v", resp, err)
	}
	return st, nil
}

// ReadSubscriptions return list of node information with the given subscription url.
func ReadSubscriptions(r runner.Runner, sub string) ([]*SubscriptionResp, error) {
	data := map[string]string{"url": sub}
	jsondata, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s not valid", sub)
	}
	resp, err := web.JSONPost(r.Host.GetURL()+"/"+"readsubscriptions", jsondata)
	if err != nil {
		return nil, fmt.Errorf("post sub: %v", err)
	}
	var cfg []*SubscriptionResp
	err = json.Unmarshal(resp, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%q not valid", resp)
	}
	return cfg, nil
}

// GetResult return the newest speed test result.
func GetResult(r runner.Runner) (*Result, error) {
	resp, err := web.Get(r.Host.GetURL() + "/" + "getresults")
	if err != nil {
		return nil, fmt.Errorf("get result: %v", err)
	}
	var result Result
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("%q not valid", resp)
	}
	return &result, nil
}

// StartTest post a speed test request with the given node config.
// Because of the blocking backend, This function is designed to be called as goroutine.
// If state is not null it will return status: "running" or "done".
// Error message response will be wrapped into a new error.
func StartTest(r runner.Runner, startCFG *StartConfigs) (string, error) {
	r.Activate()
	defer r.HangUp()

	d, err := json.Marshal(startCFG)
	if err != nil {
		return "", errors.New("invalid start config")
	}

	resp, err := web.JSONPostWithTimeout(r.Host.GetURL()+"/"+"start", d, 0)
	if err != nil {
		return "", fmt.Errorf("post speedtest start request: %v", err)
	}

	var state Status
	err = json.Unmarshal(resp, &state)
	if err != nil {
		return "", fmt.Errorf("decode %s: %v", state, err)
	}
	if state.State != "" {
		return state.State, nil
	}
	return "", errors.New(state.Error)
}

// IncludeRemarks will select all the configs with the given remarks.
func IncludeRemarks(configs []*SubscriptionResp, incRems []string) []*SubscriptionResp {
	var newcfg []*SubscriptionResp
	for _, e := range incRems {
		for _, c := range configs {
			if strings.Contains(c.Config.Remarks, e) {
				newcfg = append(newcfg, c)
				break
			}
		}
	}
	return newcfg
}

// ExcludeRemarks will select all the configs excluded the given remarks.
func ExcludeRemarks(configs []*SubscriptionResp, excludeRemarks []string) []*SubscriptionResp {
	var newcfg []*SubscriptionResp
	remarksLength := len(excludeRemarks)
	for _, c := range configs {
		for i, mark := range excludeRemarks {
			if mark == c.Config.Remarks {
				break
			}
			if i == remarksLength-1 {
				newcfg = append(newcfg, c)
			}
		}
	}
	return newcfg
}
