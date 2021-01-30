package main

import (
	"os"
)

func init() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
}

func ExampleHello() {
	cli.Main([]string{})
	cli.Main([]string{"arg"})
	cli.Main([]string{"hello"})
	cli.Main([]string{"hello", "world"})
	cli.Main([]string{"hello.world"})
	cli.Main([]string{"hello", "script"})
	cli.Main([]string{"hello", "python"})
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
	cli.Main([]string{"echo", "hello"})
	cli.Main([]string{"echo", "twice", "hi"})
	// Output:
	// hello
	// hi
	// hi
}

func ExampleEnv() {
	os.Unsetenv("MESSAGE")
	cli.Main([]string{"echo.env.message"})
	os.Setenv("MESSAGE", "system")
	cli.Main([]string{"echo.env.message"})
	os.Unsetenv("MESSAGE")
	cli.Main([]string{"echo", "env", "message2"})
	os.Setenv("MESSAGE", "system")
	cli.Main([]string{"echo", "env", "message2"})
	// Output:
	// message
	// system
	// message2
	// system
}

func ExampleAlias() {
	cli.Main([]string{"alias", "helloworld", "alias"})
	// Output:
	// hello world alias
	// bye world alias
}
