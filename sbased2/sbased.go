package sbased2

type bindState uint8

const (
	unbounded bindState = iota
	unused
	used
)

// Binder is a string based binder
type Binder struct {
	internalOpts
	values stringValues
}

type stringValues map[string]*stringValue

type stringValue struct {
	location string
	value    string
	state    bindState
}

func newStringValue(value *StringValue) *stringValue {
	return &stringValue{
		location: value.Location,
		value:    value.Value,
		state:    unbounded,
	}
}

func newStringValues(
	strValues StringValues,
	aliases map[string]string,
) stringValues {
	if len(strValues) == 0 {
		return nil
	}

	values := make(stringValues, len(strValues))

	for _, value := range strValues {
		actualKey, found := aliases[value.Key]

		if !found {
			actualKey = value.Key
		}

		values[actualKey] = newStringValue(value)
	}

	return values
}

// New creates a new string based binder
func New(provider StringValuesProvider, options ...Option) (*Binder, error) {
	if provider == nil {
		return nil, ErrNilProvider
	}

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
