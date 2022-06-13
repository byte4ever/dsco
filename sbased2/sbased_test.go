package sbased2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func nilProviderBuilder(*testing.T) StringValuesProvider {
	return nil
}

func valuesProviderBuilder(
	values StringValues,
	call bool,
) func(t *testing.T) StringValuesProvider {
	return func(t *testing.T) StringValuesProvider {
		t.Helper()

		mockedProvider := NewMockStringValuesProvider(t)

		if call {
			mockedProvider.
				On("GetStringValues").
				Return(values).
				Once()
		}

		return mockedProvider
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	type argsT struct {
		providerBuilder func(t *testing.T) StringValuesProvider
		options         []Option
	}

	tests := []struct {
		name                  string
		args                  argsT
		expectedState         *Binder
		expectedError         error
		expectedErrorContains []string
	}{
		{
			name: "success",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues{
						{
							Key:      "key1",
							Location: "loc-key1",
							Value:    "val1",
						},
						{
							Key:      "key2",
							Location: "loc-key2",
							Value:    "val2",
						},
						{
							Key:      "key3",
							Location: "loc-key3",
							Value:    "val3",
						},
					}, true,
				),
				options: nil,
			},
			expectedState: &Binder{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "val1",
					},
					"key2": {
						location: "loc-key2",
						value:    "val2",
					},
					"key3": {
						location: "loc-key3",
						value:    "val3",
					},
				},
			},
		},
		{
			name: "success no keys",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues{}, true,
				),
				options: nil,
			},
			expectedState: &Binder{
				values: stringValues(nil),
			},
		},
		{
			name: "success nil keys",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues(nil), true,
				),
				options: nil,
			},
			expectedState: &Binder{
				values: stringValues(nil),
			},
		},
		{
			name: "with aliases",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues{
						{
							Key:      "alias1",
							Location: "loc-key1",
							Value:    "val1",
						},
						{
							Key:      "key2",
							Location: "loc-key2",
							Value:    "val2",
						},
						{
							Key:      "alias3",
							Location: "loc-key3",
							Value:    "val3",
						},
					}, true,
				),
				options: []Option{
					WithAliases(
						map[string]string{
							"alias1": "key1",
							"alias3": "key3",
						},
					),
				},
			},
			expectedState: &Binder{
				values: stringValues{
					"key1": {
						location: "loc-key1",
						value:    "val1",
					},
					"key2": {
						location: "loc-key2",
						value:    "val2",
					},
					"key3": {
						location: "loc-key3",
						value:    "val3",
					},
				},
				internalOpts: internalOpts{
					aliases: map[string]string{
						"alias1": "key1",
						"alias3": "key3",
					},
				},
			},
		},
		{
			name: "with no aliases provided",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues{}, false,
				),
				options: []Option{
					WithAliases(
						map[string]string{},
					),
				},
			},
			expectedState: nil,
			expectedError: ErrNoAliasesProvided,
		},
		{
			name: "with nil aliases provided",
			args: argsT{
				providerBuilder: valuesProviderBuilder(
					StringValues{}, false,
				),
				options: []Option{
					WithAliases(nil),
				},
			},
			expectedState: nil,
			expectedError: ErrNoAliasesProvided,
		},
		{
			name: "with nil provider",
			args: argsT{
				providerBuilder: nilProviderBuilder,
			},
			expectedState: nil,
			expectedError: ErrNilProvider,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				// test invariants
				if tt.expectedState == nil && tt.expectedError == nil ||
					tt.expectedState != nil && tt.expectedError != nil {
					t.Error(
						"test invariant failure detected " +
							"cannot get both error and result",
					)
				}

				provider := tt.args.providerBuilder(t)

				binder, err := New(provider, tt.args.options...)

				if tt.expectedError != nil {
					require.ErrorIs(t, err, tt.expectedError)

					for _, contain := range tt.expectedErrorContains {
						require.ErrorContainsf(
							t,
							err,
							contain,
							"error <%v> does not contain %q",
							err,
							contain,
						)
					}

					require.Nil(t, binder)
					return
				}

				require.NoError(t, err)
				require.Equal(t, tt.expectedState, binder)
			},
		)
	}
}
