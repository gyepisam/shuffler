<A name="toc1-0" title="What" />
# What

Package shuffler shuffles items using the standard Go math/rand.Perm shuffler, which uses
the Fisher-Yates algorithm.

<A name="toc1-6" title="Why" />
# Why

Shuffled items can be anchored by position or relationship to other items.
This is particularly useful for market researchers doing surveys but may be useful to others.

<A name="toc1-12" title="How" />
# How

Here is a complete example, also found in the example directory:

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

<A name="toc1-19" title="When" />
# When

Though I have not seen other implementations of shufflers with anchoring, I imagine they've been around for a while.
I have written several implementations over the years in various languages.

<A name="toc1-25" title="Who" />
# Who

Shuffler is written by Gyepi Sam <self-github@gyepi.com>

  
