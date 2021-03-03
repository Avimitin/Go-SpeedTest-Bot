package speedtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-speedtest-bot/internal/web"
	"os"
	"path"
	"strings"
)

// Ping will test if the given host is accessible or not
func Ping(h Host) bool {
	resp, err := web.Get(path.Join(h.GetURL(), "getversion"))
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
func GetStatus(h Host) (*Status, error) {
	resp, err := web.Get(path.Join(h.GetURL(), "status"))
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
func ReadSubscriptions(h Host, sub string) ([]*SubscriptionResp, error) {
	data := map[string]string{"url": sub}
	jsondata, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s not valid", sub)
	}
	resp, err := web.JSONPost(path.Join(h.GetURL(), "readsubscriptions"), jsondata)
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
func GetResult(h Host) (*Result, error) {
	resp, err := web.Get(path.Join(h.GetURL(), "getresults"))
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
// If error happen function will pass error message into status chan.
// If state is not null it will pass status: ["running" / "done"].
// If you are calling this method without authorized of given unexpected config it will pass error message.
func StartTest(h Host, startCFG *StartConfigs, statusChan chan string) {
	d, err := json.Marshal(startCFG)
	if err != nil {
		e := sendStatus(statusChan, "invalid start config")
		if e != nil {
			fmt.Println(e)
			os.Exit(0)
		}
		return
	}
	resp, err := web.JSONPostWithTimeout(path.Join(h.GetURL(), "start"), d, 0)
	if err != nil {
		sendStatus(statusChan, fmt.Sprintf("speed test failed, response: %q", resp))
		return
	}
	var state Status
	err = json.Unmarshal(resp, &state)
	if err != nil {
		sendStatus(statusChan, fmt.Sprintf(""))
		return
	}
	if state.State != "" {
		sendStatus(statusChan, state.State)
	} else {
		sendStatus(statusChan, state.Error)
	}
}

func sendStatus(status chan string, content string) error {
	var err error
	if status == nil {
		return errors.New("nil status channel")
	}
	defer func() {
		msg := recover()
		err = fmt.Errorf("send status: %v", msg)
	}()
	status <- content
	return err
}

// IncludeRemarks will select all the configs with the given remarks.
func IncludeRemarks(configs []*SubscriptionResp, incRems []string) []*SubscriptionResp {
	var newcfg []*SubscriptionResp
	for _, e := range incRems {
		for _, c := range configs {
			if strings.Contains(c.Config.Remarks, e) {
				newcfg = append(newcfg, c)
			}
		}
	}
	return newcfg
}

// ExcludeRemarks will select all the configs excluded the given remarks.
func ExcludeRemarks(configs []*SubscriptionResp, excRem []string) []*SubscriptionResp {
	incRems := IncludeRemarks(configs, excRem)
	if incRems == nil {
		return nil
	}
	var set1 = map[string]*SubscriptionResp{}
	for _, i := range incRems {
		set1[i.Config.Remarks] = i
	}
	var set2 = map[string]*SubscriptionResp{}
	for _, c := range configs {
		set2[c.Config.Remarks] = c
	}
	if len(set1) > len(set2) {
		set1, set2 = set2, set1
	}
	var excludeCFG []*SubscriptionResp
	for k := range set1 {
		if _, ok := set2[k]; ok {
			delete(set2, k)
		}
	}
	for _, v := range set2 {
		excludeCFG = append(excludeCFG, v)
	}
	return excludeCFG
}
