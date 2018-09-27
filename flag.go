package main

import "fmt"

var flags = map[string]flag{}

type flag struct {
	Long     string
	Short    rune
	Complete compFunc
}

func insertFlag(f flag) {
	if f.Long != "" {
		flags["--"+f.Long] = f
	}
	if f.Short != 0 {
		flags[fmt.Sprintf("-%c", f.Short)] = f
	}
}
