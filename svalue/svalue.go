package svalue

// Value is a string value.
type Value struct {
	Location string
	Value    string
}

// Values is a string value collection.
type Values map[string]*Value
