# Z

![Go](https://github.com/zakuro9715/z/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/zakuro9715/z/branch/main/graph/badge.svg?token=K937ZYFF9Z)](https://codecov.io/gh/zakuro9715/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/zakuro9715/z)](https://goreportcard.com/report/github.com/zakuro9715/z)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


Z is a simple and useful task runner.

# Overview

## Philosophy

- Simple
- Easy
- Intuitive
- Useful

## Features

- Nested tasks
- Default task
- Alias
- Shorthand
- And more...

## Installation

Via gobinaries

```
curl -sSL gobinaries.com/zakuro9715/z | sh
```

By go install

```
go install github.com/zakuro9715/z
```

# Compare with other tools

## [Robo](https://github.com/tj/robo)

Z is strongly inspired by Robo

### Good

- Easy to use
- Simple configuration
- Easy to install

### Bad

- No nested tasks
- Not enough features
- No default task

## Make

### Good

- Easy to use
- Run anyware

### Cons

- Make is not task runner
- Makefile is difficult
- No nested tasks

## npm script

### Pros

- Easy to use
- Simple configuration
- No extra tool is required in nodejs project.

### Cons

- Not suitable for other than nodejs project
- No nested tasks
- script must be one-liner

## [Task](https://taskfile.dev) (go-task/task)

### Pros

- Many features
- Good documentation

### Cons

- No nested tasks
- Too many features

# Usage

```
z tasks... args...
```

## Config

### Run

Run with specified shell (default: sh)

```yaml
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

`run` can be omitted

```
tasks:
  hello: echo hello
```

### Shorthand

```yaml
tasks:
    hello.world: echo hello world
```

```
$ z hello world
hello world
$ z hello.world
hello world
```


### Args or flags

You can specify args and flags. They are passed to each commands

```
$ z hello world

# It runs
sh -c "echo hello1 world"
sh -c "echo hello2 world"
```

### Default task

You can use default task

```yaml
default: hello.world
tasks:
    hello:
        tasks:
            world: echo hello world
```

```
$ z
hello world
```

### Task Alias

```yaml
tasks:
    hello.world: echo hello world
    helloworld:
        z: hello.world
```

```
$z helloworld
hello world
```
### PATH

You can specify additional PATH

```
tasks:
    hello:
        path: ./bin
        run: command-in-bin-dir
```


# Use cases

- [zakuro9715/v-zconfig](https://github.com/zakuro9715/v-zconfig/blob/main/z.yaml)

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
default: hello.world               # Default task. hello.world -> z hello world
env:
  MESSAGE: message                 # It used if environment variable does not exist.
tasks:                             # Task list
  hello:                           # Task name
    desc: Say hello                # Task description
    run:                           # Commands to run
      - echo hello                 # `bash -c "echo hello {args}`
      - echo bye                   # `bash -c "echo bye {args}"`
    args:
      required: true               # Required one more arguments
      default: you                 # Default argument
    hooks:                         # hooks
      pre: echo saying hello       # pre hook
      post: echo said hello        # post hook
    tasks:                         # Sub task list
      script:
        run: examples/hello.sh     # Run script
      script.with_path:
        path: examples             # Add path
        run: hello.sh
      python:
        shell: python
        run: print('hello python')

  hello.world:                     # Sub task shorthand (Task will be 'z hello world')
    run:
      - z hello -- world           # Args are passed all commands
                                   # so it runs 'bash -c "echo hello world"' and 'bash -c "echo bye world"
                                   # after -- is args (not subtask name)
  echo: echo                       # Shorthand command ('run' can be omitted').
  echo.twice:                      # Multi commands can be used
    - echo
    - echo
  echo.env.message: echo $MESSAGE  # use env
  echo.env.message2:
    env: MESSAGE=message2          # task local default env
    run: echo $MESSAGE
  alias.helloworld:
    z: hello.world                 # Alias to other task
```

```examples/npm.yaml
default: npm.script
tasks:
  npm.script: npm run
```
