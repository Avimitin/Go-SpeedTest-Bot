package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlPak "net/url"
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

// JSONPost will post a JSON request with the url and data given by user. it will return error if timeout 30 seconds.
func JSONPost(url string, data []byte) ([]byte, error) {
	return JSONPostWithTimeout(url, data, 30)
}

// JSONPostWithTimeout will post a JSON request with the url and data given by user, it will return error if reach time limit.
func JSONPostWithTimeout(url string, data []byte, timeout int) ([]byte, error) {
	c := GetClient(timeout)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[WebPostError]%v", err)
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("[WebPostError]%v", err)
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[WebPostError]%v", err)
		return nil, err
	}
	return content, nil
}
