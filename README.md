# Z

![Go](https://github.com/zakuro9715/z/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/zakuro9715/z/branch/main/graph/badge.svg?token=K937ZYFF9Z)](https://codecov.io/gh/zakuro9715/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/zakuro9715/z)](https://goreportcard.com/report/github.com/zakuro9715/z)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


Simple task runner

# Installation

Via gobinaries

```
curl -sSL gobinaries.com/zakuro9715/z | sh
```

Build from source

```
go get github.com/zakuro9715/z
```

# Usage

```
z tasks... args...
```

# Config

## Run

Run with specified shell (default: sh)

```
tasks:
    hello:
        run:
            - echo hello1
            - echo hello2
```

```
$ z hello

# It runs
sh -c "echo hello1"
sh -c "echo hello2"
```

### Args or flags

You can specify args and flags. They are passed to each commands

```
$ z hello world

# It runs
sh -c "echo hello1 world"
sh -c "echo hello2 world"
```

You can pass flags and args. They will be passed to each commands



# Examples

See also [Examples Test](./examples_test.go)

```examples/cc.yaml
tasks:
  compile:
    run:
      - clang $@
    desc: Compile
    hooks:
      pre: echo Compiling
      post: echo Compiled
    tasks:
      main:
        run:
          - z -c examples/cc.yaml compile main.c
```

```examples/hello.yaml
shell: bash                        # Shell to run commands
default: hello                     # Default task
tasks:                             # Task list
  hello:                           # Task name
    desc: Say hello                # Task description
    run:                           # Commands to run
      - echo hello                 # It runs `bash -c echo hello`
      - echo bye
    hooks:                         # hooks
      pre: echo saying hello       # pre hook
      post: echo said hello        # post hook
    tasks:                         # Sub task list
      world:                       # Sub task name (Task will be 'z hello world')
        run:
          - z hello -- world       # Args are passed all commands
                                   # so it runs 'bash -c "echo hello world"' and 'bash -c "echo bye world"
                                   # after -- is args (not subtask name)
      script:
        script: examples/hello.sh  # Run script

      python:
        shell: python
        run: print('hello python') # Run shorthand
```
