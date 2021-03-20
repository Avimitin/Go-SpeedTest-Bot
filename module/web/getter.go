package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get method will make a get request with the url given by user.
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get %s: %v", url, err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response from %s: %v", url, err)
	}
	return content, nil
}
