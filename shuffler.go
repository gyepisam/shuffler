/*
Package shuffler shuffles arrays of integers.

Array entries can be free or anchored.
The former type is shuffled and the latter type is not.
An entry can be anchored by position and it retains the same position after shuffling.
An entry could also be anchored relative its previous or next entry, in which case, it
retains the same relative position to its anchor, which can be shuffled.

Terminology
"A > B" describes a list of two items where A is anchored to its successor, B
"A B <" describes the inverse, where B is anchored to its predecessor, A.
"A B . C" describes a list where B is anchored by position.


Edge Cases:

In "A > B <"  A and B are mutually anchored. This will be converted into "A > B"
In "A < B C"  and "A B C >" The endpoints are anchored to non-existent neighbors. The anchors are removed.

Chains of references are handled correctly.
*/

package shuffler

import (
	"math/rand"
	"play/anchor"
)

// A shuffler shuffles integers (generally representing positions)
// while maintaining anchored positions.
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
// anchorType indicates how the item is anchored.
func (s *shuffler) Add(slot int, anchorType anchor.Type) {
	s.items = append(s.items, slot)
	switch anchorType {
	case anchor.Position:
		s.position[slot] = true
	case anchor.ToPrevious:
		s.after[slot-1] = slot
		s.skip[slot] = true
	case anchor.ToNext:
		s.before[slot+1] = slot
		s.skip[slot] = true
	}
}

// resolve dangling references and break mutual anchors
// A < B C becomes A B C
// A B C > becomes A B C
// A > B < becomes A > B
func (s *shuffler) resolve() {

	if slot, ok := s.after[-1]; ok {
		delete(s.after, -1)
		delete(s.skip, slot)
	}

	if slot, ok := s.before[len(s.items)]; ok {
		delete(s.before, len(s.items))
		delete(s.skip, slot)
	}

	for from, to := range s.after {
		if s.before[to] == from {
			delete(s.before, to)
			delete(s.skip, from)
		}
	}
}

// Shuffle shuffles and returns the list of items added with Add, using seed
// to initialize the random number generator.
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
			out[j] = j
			continue
		}

		// skip reserved position anchored slots
		for s.position[k] {
			k++
		}

		// insert "ToNext" anchored items before item.
		// follow chains of references
		for v, ok := s.before[j]; ok; v, ok = s.before[v] {
			out[k] = v
			k++
		}

		// place the item
		out[k] = j
		k++

		// insert "ToPrevious" anchored items after item.
		// follow chains of references
		for v, ok := s.after[j]; ok; v, ok = s.after[v] {
			out[k] = v
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
