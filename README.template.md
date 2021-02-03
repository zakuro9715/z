# Z

![Go](https://github.com/zakuro9715/z/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/zakuro9715/z/branch/main/graph/badge.svg?token=K937ZYFF9Z)](https://codecov.io/gh/zakuro9715/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/zakuro9715/z)](https://goreportcard.com/report/github.com/zakuro9715/z)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)


Simple task runner

# Overview

## Philosophy

- Simple
- Easy
- Intuitive
- Useful

## Features

- Default task
- Nested tasks
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

You can pass flags and args. They will be passed to each commands

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

# Use cases

- [zakuro9715/v-zconfig](https://github.com/zakuro9715/v-zconfig/blob/main/z.yaml)

# Examples

See also [Examples Test](./examples_test.go)

__examples__
