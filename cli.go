package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/zakuro9715/nzflag"
	"github.com/zakuro9715/z/config"
	"github.com/zakuro9715/z/log"
	"github.com/zakuro9715/z/runner"
)

var exit = os.Exit

var version = "v0.4.0 or later" // Set via ldflags

const helpTextBase = `
z - Z Task runner

Usage:
  z [OPTIONS] task... [ARGS]
`

const (
	ENV_KEY_ZCONNFIG = "ZCONFIG"
	ENV_KEY_ZVERBOSE = "ZVERBOSE"
)

func isHelpFlag(s string) bool {
	return s == "-h" || s == "--help"
}

func isEnvValTrue(s string) bool {
	return len(s) > 0 && s != "0" && s != "no"
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
	fmt.Fprintf(w, "z %v\n", version)
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
	verboseFlag := isEnvValTrue(os.Getenv(ENV_KEY_ZVERBOSE))
	for ; i < len(nzargs); i++ {
		arg := nzargs[i]
		if arg.Type() != nzflag.TypeFlag {
			break
		}
		switch {
		case isHelpFlag(arg.String()):
			helpFlag = true
		case arg.Flag().Name == "V" || arg.Flag().Name == "version":
			fprintVersion(os.Stdout)
		case arg.Flag().Name == "v" || arg.Flag().Name == "verbose":
			verboseFlag = true
		case arg.Flag().Name == "c" || arg.Flag().Name == "config":
			configPath = arg.Flag().Values[0]
		default: // unknow flag
			unknownFlag = true
		}
	}

	if verboseFlag {
		log.Default.Level = log.INFO
		os.Setenv(ENV_KEY_ZVERBOSE, "1")
	}
	os.Setenv(ENV_KEY_ZCONNFIG, configPath)

	log.Info("flags:")
	log.Info("  %v", nzargs[0:i])

	config, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, yaml.FormatError(err, true, true))
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

	task := config.FindTask(config.Default)
	if i >= len(nzargs) {
		if len(config.Default) == 0 {
			fprintHelp(os.Stderr, config)
			return 1
		}
	} else {
		fullName := ""
		i2 := i
		for ; i2 < len(nzargs); i2++ {
			arg := nzargs[i2]
			if isHelpFlag(nzargs[i2].String()) {
				fprintTaskHelp(os.Stdout, task)
				return 0
			}
			if arg.String() == "--" {
				i = i2 + 1
				break
			}
			if len(fullName) > 0 {
				fullName += "."
			}
			fullName += arg.String()
			newTask := config.FindTask(fullName)
			if newTask != nil {
				i = i2 + 1
				task = newTask
			}
		}
	}

	taskArgs := make([]string, len(nzargs[i:]))
	for i, arg := range nzargs[i:] {
		taskArgs[i] = arg.String()
	}
	runner.Run(task, taskArgs)
	return 0
}
