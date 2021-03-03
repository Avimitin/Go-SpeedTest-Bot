package web

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

const (
	addr             = "http://127.0.0.1:8888"
	testAmount int32 = 2
)

var (
	sig = make(chan interface{})
)

func TestMain(m *testing.M) {
	go func() {
		m.Run()
		var done int32
		for range sig {
			if done == testAmount {
				os.Exit(0)
			}
		}
	}()
	http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(rw, "message")
		} else if r.Method == http.MethodPost {
			fmt.Fprint(rw, r.Body)
		}
	}))
}

func TestGet(t *testing.T) {
	get, err := Get(addr)
	if err != nil {
		t.Fatal(err)
	}
	want := "message"
	if string(get) != want {
		t.Errorf("get %s want %s", get, want)
	}
	sig <- 0
}

func TestPost(t *testing.T) {
	data := map[string][]string{
		"title": {"foo"},
	}
	get, err := Post(addr, data)
	if err != nil {
		t.Fatal(err)
	}
	want := `"title": "foo"`
	if string(get) != want {
		t.Errorf("get %s want %s", get, want)
	}
	sig <- 0
}
