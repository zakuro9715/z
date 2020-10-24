package main

import (
	"os"
)

func ExampleHello() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
	realMain([]string{"hello"})
	realMain([]string{"hello", "world"})
	realMain([]string{"hello", "cjk"})
	// Output:
	// hello
	// hello world
	// hello japan
	// hello korea
	// hello china
}
