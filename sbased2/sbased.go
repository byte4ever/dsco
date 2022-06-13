package sbased2

type StringValue struct {
	Key      string
	Location string
	Value    string
}

type StringValues []StringValue

type StringValuesProvider interface {
	GetStringValues() StringValues
}

type stringValue struct {
	location string
	value    string
	bounded  bool
	used     bool
}

func newStringValue(value StringValue) stringValue {
	return stringValue{
		location: value.Location,
		value:    value.Value,
		bounded:  false,
		used:     false,
	}
}

type stringValues map[string]stringValue

func newStringValues(
	values StringValues,
	aliases map[string]string,
) stringValues {
	r := make(stringValues, len(values))

	for _, value := range values {
		actualKey, found := aliases[value.Key]

		if !found {
			actualKey = value.Key
		}

		r[actualKey] = newStringValue(value)
	}

	return r
}

type Binder struct {
	internalOpts
	values stringValues
}

func New(provider StringValuesProvider, options ...Option) (*Binder, error) {
	internalOptions := internalOpts{}

	if err := internalOptions.applyOptions(options); err != nil {
		return nil, err
	}

	return &Binder{
		values: newStringValues(
			provider.GetStringValues(), internalOptions.aliases,
		),
		internalOpts: internalOptions,
	}, nil
}
