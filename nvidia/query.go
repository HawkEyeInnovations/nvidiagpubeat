package nvidia

// Map is a type alias
type Map map[string]struct{}

// Query is a struct containing the maps of settings to output from the beat
type Query struct {
	System Map
	GPU    Map
}

// NewQuery construct and returns an Query struct populated from string arrays
func NewQuery(system []string, gpu []string) Query {
	query := Query{
		System: make(Map),
		GPU:    make(Map),
	}

	for _, t := range system {
		query.System[t] = struct{}{}
	}

	for _, t := range gpu {
		query.GPU[t] = struct{}{}
	}

	return query
}
