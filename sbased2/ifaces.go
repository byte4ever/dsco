package sbased2

type StringValue struct {
	Key      string
	Location string
	Value    string
}

type StringValues []*StringValue

type StringValuesProvider interface {
	GetStringValues() StringValues
}
