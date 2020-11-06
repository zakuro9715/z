package main

import (
	"os"
)

var Version = "unset"

func main() {
	os.Exit(realMain(os.Args[1:]))
}
