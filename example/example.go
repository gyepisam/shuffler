package main

import (
	"fmt"

	"github.com/gyepisam/shuffler"

	"math/rand"
)

func main() {
	// Lightly edited list of types of shuffles, copied from Wikipedia
	shuffles := []struct {
		name   string
		anchor shuffler.Anchor
	}{
		{"Chemmy", shuffler.None}, //shuffler.None will be shuffled
		{"Corgi", shuffler.None},
		{"Faro", shuffler.None},
		{"Indian", shuffler.None},
		{"Irish", shuffler.ToPrevious},   // ToPrevious will stick to previous item
		{"Mexican", shuffler.ToPrevious}, // anchors can be chained
		{"Mongean", shuffler.None},
		{"Overhand", shuffler.Position}, //Position will keep item in spot
		{"Pile", shuffler.None},
		{"Riffle", shuffler.None},
		{"Stripping", shuffler.ToNext}, //ToNext will anchor it to next item.
		{"Wash", shuffler.None},
		{"Weave", shuffler.Position},
	}

	shuf := shuffler.New()
	for i, shuffle := range shuffles {
		shuf.Add(i, shuffle.anchor)
	}

	indices := shuf.Shuffle(rand.Int63())

	fmt.Println("Sorted list of shuffles:")
	for i, shuffle := range shuffles {
		fmt.Printf("%d %s\n", i, shuffle.name)
	}

	fmt.Println("Shuffled list of shuffles:")
	for _, j := range indices {
		fmt.Printf("%d %s\n", j, shuffles[j].name)
	}
}
