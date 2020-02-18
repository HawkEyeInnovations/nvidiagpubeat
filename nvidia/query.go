package nvidia

// Map is a type alias
type Map map[string]struct{}

// Query is a struct containing the maps of settings to output from the beat
type Query struct {
	System Map
	GPU    Map
}

// NewQuery construct and returns an empty Query struct
func NewQuery() Query {
	return Query{
		System: make(Map),
		GPU:    make(Map),
	}
}
