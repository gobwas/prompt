# prompt

[![GoDoc][godoc-image]][godoc-url]

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

For more details please visit examples folder.

[godoc-image]: https://godoc.org/github.com/gobwas/prompt?status.svg
[godoc-url]:   https://godoc.org/github.com/gobwas/prompt
