package main

import (
	"os"
	"testing"

	"github.com/zakuro9715/z/cli"
)

func init() {
	os.Setenv("ZCONFIG", "examples/hello.yaml")
}

func BenchmarkHelloExapmle(b *testing.B) {
	os.Setenv("ZSILENT", "1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Main([]string{})
		cli.Main([]string{"arg"})
		cli.Main([]string{"hello"})
		cli.Main([]string{"hello", "world"})
		cli.Main([]string{"hello.world"})
		cli.Main([]string{"hello", "script"})
		cli.Main([]string{"hello", "python"})

		os.Unsetenv("MESSAGE")
		cli.Main([]string{"echo.env.message"})
		os.Setenv("MESSAGE", "system")
		cli.Main([]string{"echo.env.message"})
		os.Unsetenv("MESSAGE")
		cli.Main([]string{"echo", "env", "message2"})
		os.Setenv("MESSAGE", "system")
		cli.Main([]string{"echo", "env", "message2"})

		cli.Main([]string{"helloworld", "alias"})
	}
}

func ExampleHello() {
	cli.Main([]string{})
	cli.Main([]string{"arg"})
	cli.Main([]string{"hello"})
	cli.Main([]string{"hello", "world"})
	cli.Main([]string{"hello.world"})
	cli.Main([]string{"hello", "script"})
	cli.Main([]string{"hello", "script", "with_path"})
	cli.Main([]string{"hello", "python"})
	// Output:
	// hello world
	// hello world arg
	// hello you
	// hello world
	// hello world
	// hello script
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

func ExampleVarAndEnv() {
	cli.Main([]string{"echo.var.value"})
	os.Unsetenv("MESSAGE")
	cli.Main([]string{"echo.env.message"})
	os.Setenv("MESSAGE", "system")
	cli.Main([]string{"echo.env.message"})
	os.Unsetenv("MESSAGE")
	cli.Main([]string{"echo", "env", "message2"})
	os.Setenv("MESSAGE", "system")
	cli.Main([]string{"echo", "env", "message2"})
	// Output:
	// value
	// message
	// system
	// message2
	// system
}

func ExampleAlias() {
	cli.Main([]string{"helloworld", "alias"})
	// Output:
	// hello world alias
}

func ExampleVerifyConfig() {
	cli.Main([]string{"--verify-config"})
	// Output:
}
