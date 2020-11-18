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
// Because of the blocking backend, This function is designed to be called as goroutine.
// If error happen function will pass error message into status chan.
// If state is not null it will pass status: ["running" / "done"].
// If you are calling this method without authorized of given unexpected config it will pass error message.
func StartTest(h *Host, startCFG *StartConfigs, status chan string) {
	d, err := json.Marshal(startCFG)
	if err != nil {
		log.Println("[JSONError]Unable to marshall data", err)
		return
	}
	fmt.Println(string(d))
	resp, err := web.JSONPostWithTimeout(h.GetURL()+"/start", d, 0)
	if err != nil {
		log.Println("[WebGetError]Unable to connect to backend")
		return
	}
	var state Status
	err = json.Unmarshal(resp, &state)
	if err != nil {
		status <- err.Error()
		return
	}
	if state.State != "" {
		status <- state.State
	} else {
		status <- state.Error
	}
}
