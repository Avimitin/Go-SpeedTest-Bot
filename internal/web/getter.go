package web

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Get method will make a get request with the url given by user.
func Get(url string) ([]byte, error) {
	if !isLegalURL(url) {
		return nil, &illegalUrl{"URL missing prefix."}
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[WebGetError]%v", err)
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[WebGetError]%v", err)
		return nil, err
	}
	return content, nil
}
