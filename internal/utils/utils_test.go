package utils

import (
	"fmt"
	"testing"
)

func TestGetIP(t *testing.T) {
	resp, err := GetIP("google.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
