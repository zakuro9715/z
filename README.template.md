# Z

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

__examples__
