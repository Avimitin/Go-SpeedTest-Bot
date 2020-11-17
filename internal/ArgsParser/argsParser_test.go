package ArgsParser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	flags := map[string]string{
		"-u": "abc",
		"-x": "hhh",
		"-m": "test",
	}
	text := "cmd -m ggg -x ddd"
	fmt.Println(Parser(flags, text))
}
