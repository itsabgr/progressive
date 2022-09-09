package span

// relation represents how two spans relationTo to each other.
type relation int

const (
	relationUnknown relation = iota

	/*
		Interval x is before Interval y.

		    +---+
		    | x |
		    +---+
		          +---+
		          | y |
		          +---+
	*/
	relationBefore

	/*
		Interval x meets Interval y.

		    +---+
		    | x |
		    +---+
		        +---+
		        | y |
		        +---+
	*/
	relationMeets

	/*
		Interval x overlaps Interval y.

		    +---+
		    | x |
		    +---+
		      +---+
		      | y |
		      +---+
	*/
	relationOverlaps

	/*
		Interval x is finished by Interval y.

		    +-----+
		    |  x  |
		    +-----+
		      +---+
		      | y |
		      +---+
	*/
	relationFinishedBy

	/*
		Interval x contains Interval y.

		    +-------+
		    |   x   |
		    +-------+
		      +---+
		      | y |
		      +---+
	*/
	relationContains

	/*
		Interval x starts Interval y.

		    +---+
		    | x |
		    +---+
		    +-----+
		    |  y  |
		    +-----+
	*/
	relationStarts

	/*
		Interval x is equal to Interval y.

		    +---+
		    | x |
		    +---+
		    +---+
		    | y |
		    +---+
	*/
	relationEqual

	/*
		Interval x is started by Interval y.

		    +-----+
		    |  x  |
		    +-----+
		    +---+
		    | y |
		    +---+
	*/
	relationStartedBy

	/*
		Interval x is during Interval y.

		      +---+
		      | x |
		      +---+
		    +-------+
		    |   y   |
		    +-------+
	*/
	relationDuring

	/*
		Interval x finishes Interval y.

		      +---+
		      | x |
		      +---+
		    +-----+
		    |  y  |
		    +-----+
	*/
	relationFinishes

	/*
		Interval x is overlapped by Interval y.

		      +---+
		      | x |
		      +---+
		    +---+
		    | y |
		    +---+
	*/
	relationOverlappedBy

	/*
		Interval x is met by Interval y.

		        +---+
		        | x |
		        +---+
		    +---+
		    | y |
		    +---+
	*/
	relationMetBy

	/*
		Interval x is after Interval y.

		          +---+
		          | x |
		          +---+
		    +---+
		    | y |
		    +---+
	*/
	relationAfter
)

//Invert a relation. Every relation has an inverse.
func (r relation) Invert() relation {
	switch r {
	case relationAfter:
		return relationBefore
	case relationBefore:
		return relationAfter
	case relationContains:
		return relationDuring
	case relationDuring:
		return relationContains
	case relationEqual:
		return relationEqual
	case relationFinishedBy:
		return relationFinishes
	case relationFinishes:
		return relationFinishedBy
	case relationMeets:
		return relationMetBy
	case relationMetBy:
		return relationMeets
	case relationOverlappedBy:
		return relationOverlaps
	case relationOverlaps:
		return relationOverlappedBy
	case relationStartedBy:
		return relationStarts
	case relationStarts:
		return relationStartedBy
	default:
		return relationUnknown
	}
}
