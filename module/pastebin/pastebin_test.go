package pastebin

import (
	"os"
	"testing"
)

func TestPaste(t *testing.T) {
	var code = "test code"
	url, err := Paste(os.Getenv("pastebin_key"), "test", &code)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}
