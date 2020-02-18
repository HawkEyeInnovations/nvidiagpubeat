package nvidia

// MockQuery returns an empty query struct
func MockQuery() Query {
	return NewQuery([]string{}, []string{})
}
