package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zakuro9715/nzflag"
	"github.com/zakuro9715/z/config"
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

func fprintTasks(w io.Writer, tasks map[string]*config.Task) {
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

func fprintHelp(w io.Writer, config *config.Config) {
	var sb strings.Builder
	sb.WriteString(helpTextBase)
	if config != nil {
		sb.WriteString("\n")
		fprintTasks(&sb, config.Tasks)
	}
	fmt.Fprintln(w, sb.String())
}

func fprintTaskHelp(w io.Writer, task *config.Task) {
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

	nzargs := (&nzflag.App{
		FlagOption: map[string]nzflag.FlagOption{
			"c":      nzflag.HasValue,
			"config": nzflag.HasValue,
		},
	}).Normalize(args)

	helpFlag := false
	unknownFlag := false
	for ; i < len(nzargs); i++ {
		arg := nzargs[i]
		if arg.Type() != nzflag.TypeFlag {
			break
		}
		switch {
		case isHelpFlag(arg.String()):
			helpFlag = true
		case arg.Flag().Name == "v" || arg.Flag().Name == "--version":
			fprintVersion(os.Stdout)
		case arg.Flag().Name == "c" || arg.Flag().Name == "config":
			configPath = arg.Flag().Values[0]
		default: // unknow flag
			unknownFlag = true
		}
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if unknownFlag {
		fprintHelp(os.Stderr, config)
		return 0
	}

	if helpFlag {
		fprintHelp(os.Stdout, config)
		return 0
	}

	taskName := config.Default
	if i >= len(nzargs) {
		if len(taskName) == 0 {
			fprintHelp(os.Stderr, config)
			return 1
		}
	} else {
		taskName = nzargs[i].String()
		i++
	}

	task, ok := config.Tasks[taskName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown task: %v\n", taskName)
		exit(1)
	}
	for ; i < len(nzargs); i++ {
		if isHelpFlag(nzargs[i].String()) {
			fprintTaskHelp(os.Stdout, task)
			return 0
		}
		if nzargs[i].String() == "--" {
			i += 1
			break
		}
		if subtask, ok := task.Tasks[nzargs[i].String()]; ok {
			task = subtask
		} else {
			break
		}
	}

	taskArgs := make([]string, len(nzargs[i:]))
	for i, arg := range nzargs[i:] {
		taskArgs[i] = arg.String()
	}
	NewTaskRunner(task).Run(taskArgs)
	return 0
}
