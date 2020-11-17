package ArgsParser

import "strings"

func Parser(flags map[string]string, s string) map[string]string {
	args := strings.Fields(s)
	if len(args) == 1 {
		return nil
	}

	for flag, _ := range flags {
		for i, arg := range args {
			if flag == arg {
				flags[flag] = args[i+1]
				args = append(args[:i], args[i+1:]...)
				break
			}
		}
	}
	return flags
}
