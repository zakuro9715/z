package main

import (
	"os"
)

func main() {
	os.Exit(realMain(os.Args[1:]))
}
