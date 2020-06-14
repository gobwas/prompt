# prompt

[![GoDoc][godoc-image]][godoc-url]
[![CI][ci-image]][ci-url]

> Ascetic and minimalistic CLI prompt library.

This is yet another CLI prompt library in Go. Its main purpose is to provide as
simple as possible API for both users of the library and users of the program
using this library.

# Usage

```go
// Simple line reading.
name, _ := prompt.ReadLine(ctx, "What's your name? ")

// Y/n confirmation.
yes, _ := prompt.Confirm(context.TODO(), name+", do you want to continue?")
if !yes {
	return
}

// Menu with single possible answer.
x, _ := prompt.SelectSingle(ctx, 
	"Which football club is your favorite?", 
	[]string{
		"FC Barcelona",
		"Spartak Moscow",
		"Manchester United",
		"Juventus",
	},
)
```

For more details please visit example folder.

# Screencast

This screencast was recorded by running binary from example folder.

<img src="https://raw.githubusercontent.com/gobwas/prompt/master/example/example.gif" width="800" style="border-radius:5px">

[godoc-image]: https://godoc.org/github.com/gobwas/prompt?status.svg
[godoc-url]:   https://godoc.org/github.com/gobwas/prompt
[ci-image]:    https://github.com/gobwas/prompt/workflows/CI/badge.svg?branch=master
[ci-url]:      https://github.com/gobwas/prompt/actions?query=workflow%3ACI
