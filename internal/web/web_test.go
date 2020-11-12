package web

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	resp, err := Get("https://baidu.com")
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	fmt.Printf("%s", string(resp))
}

func TestPost(t *testing.T) {
	data := map[string][]string{
		"title":  {"foo", "bar"},
		"userID": {"123"},
	}
	resp, err := Post("https://jsonplaceholder.typicode.com/posts", data)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%s", string(resp))
}
