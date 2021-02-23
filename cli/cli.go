package cli

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

var Version string // Set via wain.init

const helpTextBase = `
z - Z Task runner

Usage:
  z [OPTIONS] task... [ARGS]

Options:
  -c config, --config=config Specify config file [default: z.yaml, env: ZCONFIG]
  -v, --verbose              Enbale verbose log [env: ZVERBOSE]
  --silent                   Supress output [env: ZSILENT]
  -h, --help                 Print help
  -V, --verison              Print version
`

const (
	ENV_KEY_ZCONNFIG = "ZCONFIG"
	ENV_KEY_ZVERBOSE = "ZVERBOSE"
	ENV_KEY_ZSILENT  = "ZSILENT"
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
	fmt.Fprintf(w, "z %v\n", Version)
}

// exit_code == -1 means don't exit
func processFlags(nzargs []nzflag.Value) (i int, _ *config.Config, _ *runner.Config, code int) {
	configPath := "z.yaml"
	if p, ok := os.LookupEnv("ZCONFIG"); ok {
		configPath = p
	}

	helpFlag := false
	verboseFlag := isEnvValTrue(os.Getenv(ENV_KEY_ZVERBOSE))
	rconfig := &runner.Config{Silent: isEnvValTrue(os.Getenv(ENV_KEY_ZSILENT))}
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
			code = 0
			return
		case arg.Flag().Name == "v" || arg.Flag().Name == "verbose":
			verboseFlag = true
		case arg.Flag().Name == "c" || arg.Flag().Name == "config":
			configPath = arg.Flag().Values[0]
		case arg.Flag().Name == "silent":
			rconfig.Silent = true
		default: // unknow flag
			goto parse_end
		}
	}
parse_end:

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
		return i, config, rconfig, 1
	}

	if helpFlag {
		fprintHelp(os.Stdout, config)
		return i, config, rconfig, 0
	}

	return i, config, rconfig, -1
}

func findTask(i *int, nzargs []nzflag.Value, config *config.Config) (task *config.Task, code int) {
	task, _ = config.FindTask(config.Default)
	{
		fullName := ""
		i2 := *i
		for ; i2 < len(nzargs); i2++ {
			arg := nzargs[i2]
			if isHelpFlag(nzargs[i2].String()) {
				fprintTaskHelp(os.Stdout, task)
				return nil, 0
			}
			if arg.String() == "--" {
				*i = i2 + 1
				break
			}
			if len(fullName) > 0 {
				fullName += "."
			}
			fullName += arg.String()
			if newTask, err := config.FindTask(fullName); err == nil {
				*i = i2 + 1
				task = newTask
				if len(task.AliasTo) > 0 {
					task, _ = config.FindTask(task.AliasTo)
				}
			}
		}
	}
	if task == nil {
		if *i < len(nzargs) {
			fmt.Fprintf(os.Stderr, "Unknown task: %v\n\n", nzargs[*i].String())
		}
		fprintHelp(os.Stderr, config)
		return nil, 1
	}
	return task, -1
}

func Main(args []string) (code int) {
	nzargs := (&nzflag.App{
		FlagOption: map[string]nzflag.FlagOption{
			"c":      nzflag.HasValue,
			"config": nzflag.HasValue,
		},
	}).Normalize(args)

	var i int
	var task *config.Task
	var config *config.Config
	var rconfig *runner.Config
	if i, config, rconfig, code = processFlags(nzargs); code >= 0 {
		return
	}

	if task, code = findTask(&i, nzargs, config); code >= 0 {
		return
	}

	taskArgs := make([]string, len(nzargs[i:]))
	for i, arg := range nzargs[i:] {
		taskArgs[i] = arg.String()
	}
	if err := runner.Run(rconfig, task, taskArgs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
