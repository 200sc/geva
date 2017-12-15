package unique

// A Node can determine how closely related
// other nodes are to itself
type Node interface {
	// Distance reports the distance from this
	// node to another, along with whether the
	// comparison is a valid comparison-- if not,
	// the distance is meaningless.
	Distance(Node) (float64, bool)
}
