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


### Args

Args will be passed. You can use it as `$@`. If args.passthrough == true, args will be embedded in each commands

```
$ cat z.yaml
tasks:
    hello: echo hello $@
    hiho:
        run:
            - echo hi
            - echo ho
        args:
            passthrough: true
$ z hello world
# It runs
sh -c "echo hello" "sh" "world"
$ z hiho world
# It runs
sh -c "echo hi world"
sh -c "echo ho world"
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

Just use z as command. Config will be inherited.

```yaml
# myconfig.yaml
tasks:
    hello.world: echo hello world
    helloworld: z hello world # same as `z --config=myconfig.yaml hello world`
```

```
$z --config=myconfig.yaml helloworld
hello world
```

Of course, you can use another config

```yaml
tasks:
    hello:
        tasks:
            a: z --config=a.yaml hello
            b: z --config=b.yaml hello
```

### Env

```
env:
    - KEY=VALUE
tasks:
    echo: echo $KEY
```

```
$ z echo
VALUE
```

### Variable

```
var:
    seq: seq 3
tasks:
    count: {{seq}} | cat  # seq 10
```

```
$ z count
1
2
3
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

__examples__
