package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/gobwas/prompt"
)

func main() {
	ctx := context.TODO()

	// Simple line reading.
	name, _ := prompt.ReadLine(ctx, "What's your name? ")

	// Y/n confirmation.
	yes, _ := prompt.Confirm(context.TODO(), name+", do you want to continue?")
	if !yes {
		return
	}

	clubs := []string{
		"FC Barcelona",
		"Spartak Moscow",
		"Manchester United",
		"Juventus",
	}
	// Menu with single possible answer.
	x, _ := prompt.SelectSingle(ctx, "Which football club is your favorite?", clubs)

	// Menu with multiple possible answers.
	xs, _ := prompt.SelectMultiple(ctx, "Which won Champions League?", clubs)

	if !contains(xs, x) {
		fmt.Println(clubs[x] + " haven't won Champions League yet :(")
	} else {
		fmt.Println(clubs[x] + " have won Champions League!")
	}
}

func contains(xs []int, x int) bool {
	sort.Ints(xs)
	i := sort.Search(len(xs), func(i int) bool {
		return xs[i] >= x
	})
	return i < len(xs) && xs[i] == x
}
