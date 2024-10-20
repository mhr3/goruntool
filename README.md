# goruntool

goruntool is a utility designed to facilitate the use of Go tools that require a `replace` directive
in their go.mod file. This tool is particularly useful when used with `go generate`, as it allows you
to run Go tools that `go run` does not directly support.

## Usage

To use goruntool, you can include it in your `go generate` directives. Here is an example:

```go
//go:generate go run github.com/mhr3/goruntool@v0.1.0 github.com/mhr3/gocc/cmd/gocc@v0.11.0 --arg1 --arg2
```

In this example, goruntool is used to run the `gocc` tool with specific arguments. Running `go generate`
will build the specified tool and invoke it with the provided arguments.

## Why Use goruntool?

Without goruntool, you might encounter the following error message when running `go generate`
with tools that use a `replace` directive in their go.mod file:

```
The go.mod file for the module providing named packages contains one or
more replace directives. It must not contain directives that would cause
it to be interpreted differently than if it were the main module.
```

This issue is discussed in detail in [golang/go#44840](https://github.com/golang/go/issues/44840).

## How It Works

goruntool performs the following steps:

1. Parses the module and version from the provided arguments.
2. Clones the specified repository at the given version into a temporary directory (requires git).
3. Builds the module.
4. Runs the built binary with the provided arguments.
