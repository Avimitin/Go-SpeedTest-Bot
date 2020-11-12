package web

import (
	"net/http"
	"time"
)

type illegalUrl struct {
	msg string
}

func (i *illegalUrl) Error() string {
	return i.msg
}

func isLegalURL(url string) bool {
	return url[0:7] == "http://" || url[0:8] == "https://"
}

// GetClient return a client with given parameter.
// Require an integer value as a timeout to calculate seconds.
func GetClient(timeout int) *http.Client {
	to := time.Duration(timeout) * time.Second
	return &http.Client{
		Timeout: to,
	}
}
