package main

import (
	"os"
)

func ExampleHello() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
	realMain([]string{})
	realMain([]string{"hello"})
	realMain([]string{"hello", "world"})
	realMain([]string{"hello.world"})
	realMain([]string{"hello", "script"})
	realMain([]string{"hello", "python"})
	// Output:
	// hello you
	// bye you
	// hello you
	// bye you
	// hello world
	// bye world
	// hello world
	// bye world
	// hello script
	// hello python
}
