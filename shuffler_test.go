package shuffler

import (
	"fmt"
	"math/rand"
	"play/anchor"
	"regexp"
	"strings"
	"testing"
)

// Shuffle based on specification string.
func shuffleSpec(in string, shuf *shuffler) (string, error) {

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
			top.anchorType = anchor.Position
		case '>':
			top.anchorType = anchor.ToNext
		case '<':
			top.anchorType = anchor.ToPrevious
		default:
			top = &choice{string(r), anchor.None}
			choices = append(choices, top)
		}
	}

	if shuf == nil {
		shuf = New()
	}
	
	for j, choice := range choices {
		shuf.Add(j, choice.anchorType)
	}
	shuffled := shuf.Shuffle(rand.Int63())

	tmp := make([]string, len(shuffled))
	for j, k := range shuffled {
		tmp[j] = choices[k].text
	}

	return strings.Join(tmp, ""), nil
}


type choice struct {
	text       string
	anchorType anchor.Type
}

type data struct {
	name string
	in   string //input
	re   string //output regex
}

var anchored = []data{
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

func TestAnchored(t *testing.T) {
	for i, v := range anchored {

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



var Cases = [][]choice{
	// no choices
	{},

	// one choice, not anchored
	{{"blue", anchor.None}},

	// one choice, anchored
	{{"banana", anchor.Position}},

	// two choices, neither anchored
	{{"male", anchor.None}, {"female", anchor.None}},

	// two choices, last anchored
	{{"male", anchor.None}, {"female", anchor.Position}},

	// two choices, first anchored
	{{"male", anchor.Position}, {"female", anchor.None}},

	// two choices, both anchored
	{{"male", anchor.Position}, {"female", anchor.Position}},

	// three choices, middle one anchored
	{{"coke", anchor.None}, {"water", anchor.Position}, {"tea", anchor.None}},

	// four choices, none anchored
	{{"internet", anchor.None}, {"television", anchor.None}, {"radio", anchor.None}, {"newspapers", anchor.None}},

	// five choices, last one anchored
	{{"mazda", anchor.None}, {"toyota", anchor.None}, {"miata", anchor.None}, {"ford", anchor.None}, {"none of the above", anchor.Position}},
}

func TestShuffle(t *testing.T) {
	seed := rand.Int63()

	for k, list := range Cases {
		shuf := New()
		items := make([]int, len(list))
		anchored := make([]int, 0)
		for i, choice := range list {
			if choice.anchorType == anchor.Position {
				anchored = append(anchored, i)
			}
			shuf.Add(i, choice.anchorType)
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
			t.Logf("%d: %s %s\n", i, choice.text, choice.anchorType)
		}
		t.Log("\nOutput\n")
		for _, j := range out {
			choice := list[j]
			t.Logf("%d: %s %s\n", j, choice.text, choice.anchorType)
		}
	}
}

var spec = strings.Repeat("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
var N = 1000

func BenchmarkNew(b *testing.B) {
	for i := 0; i < N; i++ {
		shuffleSpec(spec, New())
	}
}
