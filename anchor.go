package shuffler

// Anchor determines whether an item can be shuffled or it's behaviour otherwise.
type Anchor uint

// Anchor types
const (
	// Not anchored. Will be shuffled
	None Anchor = iota

	// Anchored in position. Will remain in same position after shuffling.
	Position

	// Anchored to Previous item. Always succeed previous item, even if shuffled.
	ToPrevious

	// Anchored to Next item. Alway precede next item, even if shuffled.
	ToNext
)

// String representation of anchor.
func (t Anchor) String() string {
    switch(t) {
        case None:
          return "none"
        case Position:
          return "position"
        case ToPrevious:
          return "to previous"
        case ToNext:
          return "to next"
    }
    return "unknown"
}
