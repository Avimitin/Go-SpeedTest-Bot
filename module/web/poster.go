package web

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	urlPak "net/url"
	"time"
)

// Post will make a post request with the url and data given by user
func Post(url string, d urlPak.Values) ([]byte, error) {
	resp, err := http.PostForm(url, d)
	if err != nil {
		return nil, fmt.Errorf("post %s: %v", url, err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response from %s: %v", url, err)
	}
	return content, nil
}

// JSONPost will post a JSON request with the url and data given by user.
// it will return error if timeout 30 seconds.
func JSONPost(url string, data []byte) ([]byte, error) {
	return JSONPostWithTimeout(url, data, 30)
}

// JSONPostWithTimeout will post a JSON request with the url and data given by user,
// it will return error if reach time limit.
func JSONPostWithTimeout(url string, data []byte, timeout time.Duration) ([]byte, error) {
	var req *http.Request
	var err error
	if timeout == 0 {
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	}
	if err != nil {
		return nil, fmt.Errorf("new request %s: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request %s: %v", url, err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response from %s: %v", url, err)
	}

	return content, nil
}
