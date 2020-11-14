package web

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	urlPak "net/url"
	"strings"
)

// Post will make a post request with the url and data given by user
func Post(url string, data map[string][]string) ([]byte, error) {
	c := GetClient(30)
	d := urlPak.Values{}
	for k, v := range data {
		d[k] = v
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(d.Encode()))
	if err != nil {
		log.Printf("[WebPostError]%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

// JSONPost will post a JSON request with the url and data given by user
func JSONPost(url string, data []byte) ([]byte, error) {
	c := GetClient(30)
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
