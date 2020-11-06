package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var exit = os.Exit

var version = "unset" // Set via ldflags

const helpTextBase = `
z - Z Task runner

Usage:
  z [OPTIONS] task... [ARGS]
`

func isHelpFlag(s string) bool {
	return s == "-h" || s == "--help"
}

func showTextAndExit(code int, text string) {
	if code == 0 {
		fmt.Println(text)
	} else {
		fmt.Fprintln(os.Stderr, text)
	}
	exit(code)
}

func fprintTasks(w io.Writer, tasks map[string]*Task) {
	fmt.Fprintln(w, "Tasks:")
	maxNameLen := 0
	for k := range tasks {
		if len(k) > maxNameLen {
			maxNameLen = len(k)
		}
	}
	for k, t := range tasks {
		fmt.Fprintf(w,
			"  %v%v - %v\n",
			k,
			strings.Repeat(" ", maxNameLen-len(k)),
			t.Description,
		)
	}
}

func fprintHelp(w io.Writer, config *Config) {
	var sb strings.Builder
	sb.WriteString(helpTextBase)
	if config != nil {
		sb.WriteString("\n")
		fprintTasks(&sb, config.Tasks)
	}
	fmt.Fprintln(w, sb.String())
}

func fprintTaskHelp(w io.Writer, task *Task) {
	var sb strings.Builder
	sb.WriteString(task.Description)
	if len(task.Tasks) > 0 {
		fprintTasks(&sb, task.Tasks)
	}
	fmt.Fprint(w, sb.String())
}

func fprintVersion(w io.Writer) {
	fmt.Fprintf(w, "z task runner %v\n", version)
}

func realMain(args []string) int {
	i := 0
	configPath := "z.yaml"
	if p, ok := os.LookupEnv("ZCONFIG"); ok {
		configPath = p
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	for ; i < len(args); i++ {
		arg := args[i]
		if arg[0] != '-' {
			break
		}
		switch {
		case isHelpFlag(arg):
			fprintHelp(os.Stdout, config)
			return 0
		case arg == "-v" || arg == "--version":
			fprintVersion(os.Stdout)
		default: // unknow flag
			fprintHelp(os.Stderr, nil)
			return 1
		}
	}

	if i >= len(args) {
		fprintHelp(os.Stderr, config)
		return 1
	}

	task, ok := config.Tasks[args[i]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown task: %v\n", args[i])
		exit(1)
	}
	i++
	for ; i < len(args); i++ {
		if isHelpFlag(args[i]) {
			fprintTaskHelp(os.Stdout, task)
			return 0
		}
		if args[i] == "--" {
			i += 1
			break
		}
		if subtask, ok := task.Tasks[args[i]]; ok {
			task = subtask
		} else {
			break
		}
	}
	task.Run(args[i:])
	return 0
}
