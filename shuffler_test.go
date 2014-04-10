package shuffler

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

type choice struct {
	text   string
	anchor Anchor
}

type data struct {
	name string
	in   string //input
	re   string //output regex
}

var specTests = []data{
	{"empty", "", ""},
	{"one choice, not anchored", "A", "A"},
	{"one choice, anchored", "A.", "A"},
	{"two choices, neither anchored", "AB", "AB|BA"},
	{"two choices, last anchored", "AB.", "AB"},
	{"two choices, first anchored", "A.B", "AB"},
	{"two choices, both anchored", "A.B.", "AB"},
	{"three choices, middle anchored", "AB.C", "ABC|CBA"},
	{"three choices, none anchored", "ABC", "ABC|BAC|BCA|CBA|CAB|ACB"},
	{"three choices, first anchored", "A.BC", "A.."},
	{"three choices, middle anchore", "AB.C", ".B."},
	{"three choices, last anchored", "ABC.", "..C"},
	{"Five choices, last anchored", "ABDCE.", "[^E]{4}E"},
	{"hold CD", "ABC.D.E", "[^CD]*CD[^CD]*"},
	{"Dangling First", "A<B", "AB|BA"},
	{"Dangling Last", "AB>", "AB|BA"},
	{"Mutual anchor", "A>B<", "AB"},
}

// Shuffle based on specification string.
func shuffleSpec(in string, shuf *shuffler) (string, error) {

	// remove spaces. They are for humans.
	s := strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, in)

	choices := make([]*choice, 0, len(s))
	var top *choice

	for _, r := range s {

		switch r {
		case '.', '<', '>':
			if top == nil {
				return "", fmt.Errorf("incorrect input string: %s", s)
			}
		}

		switch r {
		case '.':
			top.anchor = Position
		case '>':
			top.anchor = ToNext
		case '<':
			top.anchor = ToPrevious
		default:
			top = &choice{string(r), None}
			choices = append(choices, top)
		}
	}

	if shuf == nil {
		shuf = New()
	}

	for j, choice := range choices {
		shuf.Add(j, choice.anchor)
	}

	seed, err := Seed()
	if err != nil {
		return "", err
	}
	shuffled := shuf.Shuffle(seed)

	tmp := make([]string, len(shuffled))
	for j, k := range shuffled {
		tmp[j] = choices[k].text
	}

	return strings.Join(tmp, ""), nil
}

func TestAnchored(t *testing.T) {
	for i, v := range specTests {

		out, err := shuffleSpec(v.in, New())
		if err != nil {
			t.Fatalf("%d: Error with input %s: %s", i, v.in, err)
		}

		t.Logf("%s -> %s\n", v.in, out)

		matched, err := regexp.MatchString(v.re, out)
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Errorf("NO MATCH. Want [%s] to match re [%s]\n", out, v.re)
		}
	}
}

var shuffleTests = [][]choice{
	// no choices
	{},

	// one choice, not anchored
	{{"blue", None}},

	// one choice, anchored
	{{"banana", Position}},

	// two choices, neither anchored
	{{"male", None}, {"female", None}},

	// two choices, last anchored
	{{"male", None}, {"female", Position}},

	// two choices, first anchored
	{{"male", Position}, {"female", None}},

	// two choices, both anchored
	{{"male", Position}, {"female", Position}},

	// three choices, middle one anchored
	{{"coke", None}, {"water", Position}, {"tea", None}},

	// four choices, none anchored
	{{"internet", None}, {"television", None}, {"radio", None}, {"newspapers", None}},

	// five choices, last one anchored
	{{"mazda", None}, {"toyota", None}, {"miata", None}, {"ford", None}, {"none of the above", Position}},
}

func TestShuffle(t *testing.T) {
	seed, err := Seed()
	if err != nil {
		t.Fatal(err)
	}

	for k, list := range shuffleTests {
		shuf := New()
		items := make([]int, len(list))
		anchored := make([]int, 0)
		for i, choice := range list {
			if choice.anchor == Position {
				anchored = append(anchored, i)
			}
			shuf.Add(i, choice.anchor)
			items[i] = i

		}

		out := shuf.Shuffle(seed)

		for _, j := range anchored {
			if out[j] != items[j] {
				t.Errorf("%d: Failed. Expected matching anchor in pos %d, got %d", k, items[j], out[j])
			}
		}

		t.Log("\nInput\n")
		for i, choice := range list {
			t.Logf("%d: %s %s\n", i, choice.text, choice.anchor)
		}
		t.Log("\nOutput\n")
		for _, j := range out {
			choice := list[j]
			t.Logf("%d: %s %s\n", j, choice.text, choice.anchor)
		}
	}
}

// Ensure that shuffler is not confusing indices with items.
// by shuffling list of items that are all larger than largest index.
func TestIndexItem(t *testing.T) {

	const MAX = 100
	cache := make(map[int]bool)
	shuf := New()
	for i := 0; i < MAX; i++ {
		value := int(rand.Int31()) + MAX // all numbers must be greater than size of array
		cache[value] = true
		shuf.Add(value, None)
	}

	out := shuf.Shuffle(rand.Int63())

	if len(cache) != len(out) {
		t.Errorf("Shuffle returned %d items, want %d", len(out), len(cache))
	}

	for _, j := range out {
		if !cache[j] {
			if j < MAX {
				t.Errorf("Shuffle returned an index, not an item: %d", j)
			} else {
				t.Errorf("Shuffle returned unknown value: %d", j)
			}
		}
	}
}

// This was used to compare current constructor with a pre-allocating one.
// No difference was found so the new one was removed.
var spec = strings.Repeat("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
var N = 1000

func BenchmarkNew(b *testing.B) {
	for i := 0; i < N; i++ {
		shuffleSpec(spec, New())
	}
}
