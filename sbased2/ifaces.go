package sbased2

// StringValue is a string value.
type StringValue struct {
	Key      string
	Location string
	Value    string
}

// StringValues is a string value collection.
type StringValues []*StringValue

// StringValuesProvider defines the behaviour if a string value provider.
type StringValuesProvider interface {
	GetStringValues() StringValues
}
