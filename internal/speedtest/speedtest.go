package speedtest

import (
	"encoding/json"
	"fmt"
	"go-speedtest-bot/internal/web"
	"log"
)

// GetStatus is used for fetching backend status
func GetStatus(h *Host) (*Status, error) {
	resp, err := web.Get(h.GetURL() + "/status")
	if err != nil {
		log.Printf("[WebGetError]Unable to connect to backend")
		return nil, err
	}
	return &Status{State: string(resp)}, nil
}

// ReadSubscriptions return list of node information with the given subscription url.
func ReadSubscriptions(h *Host, sub string) ([]*SubscriptionResp, error) {
	data := map[string]string{"url": sub}
	d, err := json.Marshal(data)
	if err != nil {
		log.Println("[JSONMarshallError]", err)
		return nil, err
	}
	resp, err := web.JSONPost(h.GetURL()+"/readsubscriptions", d)
	if err != nil {
		log.Printf("[WebPostError]Unable to connect to backend, %v", err)
		return nil, err
	}
	var cfg []*SubscriptionResp
	err = json.Unmarshal(resp, &cfg)
	if err != nil {
		log.Println("[JSONError]Unable to unmarshall data", err)
		return nil, err
	}
	return cfg, nil
}

// GetResult return the newest speed test result.
func GetResult(h *Host) (*Result, error) {
	resp, err := web.Get(h.GetURL() + "/getresults")
	if err != nil {
		log.Println("[WebGetError]Unable to connect to backend")
		return nil, err
	}
	var result Result
	err = json.Unmarshal(resp, &result)
	if err != nil {
		log.Println("[JSONError]Unable to unmarshall data", err)
		return nil, err
	}
	return &result, nil
}

// StartTest post a speed test request with the given node config.
// This method won't return any response, before start a test user should call GetStatus method to fetch newest status.
func StartTest(h *Host, startCFG *StartConfigs) {
	d, err := json.Marshal(startCFG)
	if err != nil {
		log.Println("[JSONError]Unable to marshall data", err)
		return
	}
	fmt.Println(string(d))
	_, err = web.JSONPost(h.GetURL()+"/start", d)
	if err != nil {
		log.Println("[WebGetError]Unable to connect to backend")
		return
	}
}
