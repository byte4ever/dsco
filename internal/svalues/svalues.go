package svalues

// StringValue is a string value.
type StringValue struct {
	Location string
	Value    string
}

// StringValues is a string value collection.
type StringValues map[string]*StringValue

// StringValuesProvider defines the behaviour if a string value provider.
type StringValuesProvider interface {
	GetStringValues() StringValues
}
