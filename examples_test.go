package main

import (
	"os"
)

func ExampleHello() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
	realMain([]string{"hello"})
	realMain([]string{"hello", "world"})
	// Output:
	// hello
	// bye
	// hello world
	// bye world
}
