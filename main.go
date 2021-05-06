package main

import (
	"os"

	"github.com/zakuro9715/z/cli"
)

var version = "v0.8.0 or later" // Set via ldflags

func init() {
	cli.Version = version
}

func main() {
	os.Exit(cli.Main(os.Args[1:]))
}
