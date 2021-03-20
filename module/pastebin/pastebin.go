package pastebin

import (
	"fmt"
	"go-speedtest-bot/module/web"
	"net/url"
	"strings"
)

const (
	address = "https://pastebin.com/api/api_post.php"
)

// Paste creating a new paste at https://pastebin.com
// return the paste page if the request is valid.
func Paste(key string, name string, text *string) (string, error) {
	var data = url.Values{
		"api_dev_key":    {key},
		"api_paste_code": {*text},
		"api_option":     {"paste"},
		"api_paste_name": {name},
	}
	return request(data)
}

// PasteWithExpiry creating a new paste that with expiry
// at https://pastebin.com. Return the paste page if the request is valid.
func PasteWithExpiry(key string, name string, text *string, expiry string) (string, error) {
	var data = url.Values{
		"api_dev_key":           {key},
		"api_paste_code":        {*text},
		"api_option":            {"paste"},
		"api_paste_name":        {name},
		"api_paste_expire_date": {expiry},
	}
	return request(data)
}

func request(d url.Values) (string, error) {
	resp, err := web.Post(address, d)
	if err != nil {
		return "", fmt.Errorf("post to pastebin: %v", err)
	}

	var respStr = string(resp)
	if strings.Contains(respStr, "Bad API request") {
		return "", fmt.Errorf(strings.Split(string(resp), ", ")[1])
	}

	return respStr, nil
}
