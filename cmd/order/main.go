package main

import (
	"github.com/moeryomenko/squad"
)

func main() {
	group, err := squad.New(squad.WithSignalHandler())
	_, _ = group, err
}
