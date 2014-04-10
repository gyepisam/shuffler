/*
Package shuffler shuffles arrays of integers with anchoring.

Array entries can be free or anchored.
The former type is shuffled and the latter type is not.
An entry anchored by position retains the same position after shuffling.
An entry anchored relative its previous or next entry retains the same relative position
to the anchor, which can be shuffled.
It is possible to create multiple chains of anchors.

Terminology
"A > B" describes a list of two items where A is anchored to its successor, B.
"A B <" describes the inverse, where B is anchored to its predecessor, A.
"A B . C" describes a list where B is anchored by position.


Edge Cases:

In "A > B <"  A and B are mutually anchored. This will be converted into "A > B".
In "A < B C"  and "A B C >" The endpoints are anchored to non-existent neighbors. The anchors are removed.

Chains of references are handled correctly.
*/
package shuffler

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
)

// A shuffler shuffles integers (generally representing positions) while maintaining anchored positions.
type shuffler struct {
	items    []int
	position map[int]bool
	before   map[int]int
	after    map[int]int
	skip     map[int]bool
}

// New returns an initialized shuffler.
func New() *shuffler {
	return &shuffler{
		items:    make([]int, 0),
		position: make(map[int]bool),
		before:   make(map[int]int),
		after:    make(map[int]int),
		skip:     make(map[int]bool),
	}
}

// Add appends a new items to be shuffled.
// The anchor argument specifies how the item is anchored.
func (s *shuffler) Add(slot int, anchor Anchor) {
	index := len(s.items)
	s.items = append(s.items, slot)
	switch anchor {
	case Position:
		s.position[index] = true
	case ToPrevious:
		s.after[index-1] = index
		s.skip[index] = true
	case ToNext:
		s.before[index+1] = index
		s.skip[index] = true
	}
}

// resolve dangling references and break mutual anchors
// A < B C becomes A B C
// A B C > becomes A B C
// A > B < becomes A > B
func (s *shuffler) resolve() {

	if index, ok := s.after[-1]; ok {
		delete(s.after, -1)
		delete(s.skip, index)
	}

	if index, ok := s.before[len(s.items)]; ok {
		delete(s.before, len(s.items))
		delete(s.skip, index)
	}

	for from, to := range s.after {
		if s.before[to] == from {
			delete(s.before, to)
			delete(s.skip, from)
		}
	}
}

// Shuffle shuffles and returns the list of items added with Add,
// using seed to initialize the random number generator.
func (s *shuffler) Shuffle(seed int64) []int {
	l := len(s.items)
	r := rand.New(rand.NewSource(seed))
	p := r.Perm(l)
	out := make([]int, l)
	k := 0

	s.resolve()

	for _, j := range p {

		// previous and next anchored items are established by their anchors.
		if s.skip[j] {
			continue
		}

		// position anchored items stay in place.
		if s.position[j] {
			out[j] = s.items[j]
			continue
		}

		// skip reserved position anchored slots
		for s.position[k] {
			k++
		}

		// insert "ToNext" anchored items before item.
		// follow chains of references
		for v, ok := s.before[j]; ok; v, ok = s.before[v] {
			out[k] = s.items[v]
			k++
		}

		// place the item
		out[k] = s.items[j]
		k++

		// insert "ToPrevious" anchored items after item.
		// follow chains of references
		for v, ok := s.after[j]; ok; v, ok = s.after[v] {
			out[k] = s.items[v]
			k++
		}
	}

	return out
}

// old shuffle version. only handles position anchors, specified by an array.
func shuffle(seed int64, items []int, anchored []int) []int {

	anchorMap := make(map[int]bool)
	for _, j := range anchored {
		anchorMap[j] = true
	}

	l := len(items)
	out := make([]int, l)
	k := 0
	r := rand.New(rand.NewSource(seed))
	p := r.Perm(l)
	for _, j := range p {
		if anchorMap[j] {
			out[j] = j
			continue
		}
		for anchorMap[k] {
			k++
		}

		out[k] = j
		k++
	}

	return out
}

// Seed, a convenience function, produces a 64bit int value read from crypto/rand.Reader.
func Seed() (int64, error) {
	val, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, err
	}
	return val.Int64(), nil
}
