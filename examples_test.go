package main

import (
	"os"
)

func init() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
}

func ExampleHello() {
	realMain([]string{})
	realMain([]string{"arg"})
	realMain([]string{"hello"})
	realMain([]string{"hello", "world"})
	realMain([]string{"hello.world"})
	realMain([]string{"hello", "script"})
	realMain([]string{"hello", "python"})
	// Output:
	// hello world
	// bye world
	// hello world arg
	// bye world arg
	// hello you
	// bye you
	// hello world
	// bye world
	// hello world
	// bye world
	// hello script
	// hello python
}

func ExampleEcho() {
	realMain([]string{"echo", "hello"})
	realMain([]string{"echo", "twice", "hi"})
	// Output:
	// hello
	// hi
	// hi
}

func ExampleEnv() {
	os.Unsetenv("MESSAGE")
	realMain([]string{"echo.env.message"})
	os.Setenv("MESSAGE", "system")
	realMain([]string{"echo.env.message"})
	os.Unsetenv("MESSAGE")
	realMain([]string{"echo", "env", "message2"})
	os.Setenv("MESSAGE", "system")
	realMain([]string{"echo", "env", "message2"})
	// Output:
	// message
	// system
	// message2
	// system
}

func ExampleAlias() {
	realMain([]string{"alias", "helloworld", "alias"})
	// Output:
	// hello world alias
	// bye world alias
}
