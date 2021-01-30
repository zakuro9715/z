package main

import (
	"os"

	"github.com/zakuro9715/z/cli"
)

func main() {
	os.Exit(cli.Main(os.Args[1:]))
}
