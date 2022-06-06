package sbased

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

func TestProvide(t *testing.T) {
	t.Run(
		"success", func(t *testing.T) {
			sep := NewMockStrEntriesProvider(t)
			sep.On("GetEntries").Once().
				Return(
					Entries{
						"k1": &Entry{
							ExternalKey: "extK1",
							Value:       "val1",
						},
						"k2": &Entry{
							ExternalKey: "extK2",
							Value:       "val2",
						},
					},
				)

			p, err := NewBinder(sep)
			require.NoError(t, err)
			require.NotNil(t, p)

			require.Equal(
				t, p, &Binder{
					internalOpts: internalOpts{
						aliases: nil,
					},
					entries: entries{
						"k1": &entry{
							Entry: Entry{
								ExternalKey: "extK1",
								Value:       "val1",
							},
							bounded: false,
							used:    false,
						},
						"k2": &entry{
							Entry: Entry{
								ExternalKey: "extK2",
								Value:       "val2",
							},
							bounded: false,
							used:    false,
						},
					},
					provider: sep,
				},
			)
		},
	)

	t.Run(
		"success with alias", func(t *testing.T) {
			sep := NewMockStrEntriesProvider(t)
			sep.On("GetEntries").Once().
				Return(
					Entries{
						"alias1": &Entry{
							ExternalKey: "extK1Aliased",
							Value:       "val1",
						},
						"k2": &Entry{
							ExternalKey: "extK2",
							Value:       "val2",
						},
						"alias3": &Entry{
							ExternalKey: "extK3Aliased",
							Value:       "val3",
						},
					},
				)

			p, err := NewBinder(
				sep,
				WithAliases(
					map[string]string{
						"alias1": "k1",
						"alias3": "k3",
					},
				),
			)
			require.NoError(t, err)
			require.NotNil(t, p)

			require.Equal(
				t, p, &Binder{
					internalOpts: internalOpts{
						aliases: map[string]string{
							"alias1": "k1",
							"alias3": "k3",
						},
					},
					entries: entries{
						"k1": &entry{
							Entry: Entry{
								ExternalKey: "extK1Aliased",
								Value:       "val1",
							},
							bounded: false,
							used:    false,
						},
						"k2": &entry{
							Entry: Entry{
								ExternalKey: "extK2",
								Value:       "val2",
							},
							bounded: false,
							used:    false,
						},
						"k3": &entry{
							Entry: Entry{
								ExternalKey: "extK3Aliased",
								Value:       "val3",
							},
							bounded: false,
							used:    false,
						},
					},
					provider: sep,
				},
			)
		},
	)

	t.Run(
		"success with alias", func(t *testing.T) {
			option := NewMockOption(t)
			option.On("apply", mock.Anything).Once().
				Return(mockedError)

			p, err := NewBinder(
				nil,
				option,
			)
			require.ErrorIs(t, err, mockedError)
			require.Nil(t, p)
		},
	)
}

var mockedError = errors.New("mocked error")

func getBinder() *Binder {
	return &Binder{
		internalOpts: internalOpts{
			aliases: map[string]string{
				"alias1": "k1_int",
				"alias3": "k3_float64",
			},
		},
		entries: entries{
			"k1_int": &entry{
				Entry: Entry{
					ExternalKey: "extK1Aliased",
					Value:       "1234",
				},
				bounded: false,
				used:    false,
			},
			"k2_string": &entry{
				Entry: Entry{
					ExternalKey: "extK2",
					Value:       "val2",
				},
				bounded: false,
				used:    false,
			},
			"k3_float64": &entry{
				Entry: Entry{
					ExternalKey: "extK3Aliased",
					Value:       "val3",
				},
				bounded: false,
				used:    false,
			},
			"k4_slice_int": &entry{
				Entry: Entry{
					ExternalKey: "extK4",
					Value:       `[1,2,3,4,5]`,
				},
				bounded: false,
				used:    false,
			},
		},
	}
}

func TestBinder_Bind(t *testing.T) {
	mockedOriginName := "mocked"
	mockedOrigin := dsco.Origin(mockedOriginName)

	t.Run(
		"success and set", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, outVal, err := b.Bind("k1_int", true, dstValue)
			require.NoError(t, err)
			require.True(t, succeed)
			require.Equal(t, "extK1Aliased", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := outVal.Interface()
			require.IsType(t, rs, di)
			require.Equal(t, 1234, *(di.(*int)))
			require.True(t, b.entries["k1_int"].bounded)
			require.True(t, b.entries["k1_int"].used)
		},
	)

	t.Run(
		"success and don't set", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, outVal, err := b.Bind("k1_int", false, dstValue)
			require.NoError(t, err)
			require.False(t, succeed)
			require.Equal(t, "extK1Aliased", keyOut)
			require.Equal(t, mockedOrigin, origin)

			require.Equal(t, reflect.Value{}, outVal)
			require.True(t, b.entries["k1_int"].bounded)
			require.False(t, b.entries["k1_int"].used)
		},
	)

	t.Run(
		"parse error", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstValue := reflect.ValueOf(rs)

			key := "k2_string"
			origin, keyOut, succeed, outVal, err := b.Bind(key, true, dstValue)
			require.ErrorIs(t, err, ErrParse)
			require.Equal(t, reflect.Value{}, outVal)
			require.Equal(t, mockedOrigin, origin)
			require.False(t, succeed)
			require.Equal(t, "extK2", keyOut)
			require.ErrorContains(t, err, mockedOriginName)
			require.ErrorContains(t, err, "extK2")
			require.True(t, b.entries[key].bounded)
			require.False(t, b.entries[key].used)
		},
	)
	t.Run(
		"alias collision", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstValue := reflect.ValueOf(rs)

			key := "alias1"
			origin, keyOut, succeed, outVal, err := b.Bind(key, true, dstValue)
			require.ErrorIs(t, err, ErrAliasCollision)
			require.Equal(t, reflect.Value{}, outVal)
			require.Equal(t, mockedOrigin, origin)
			require.False(t, succeed)
			require.Equal(t, "", keyOut)
			require.ErrorContains(t, err, mockedOriginName)
			require.ErrorContains(t, err, key)
		},
	)

	t.Run(
		"key not found", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, outVal, err := b.Bind("not_found", true, dstValue)
			require.NoError(t, err)
			require.Equal(t, reflect.Value{}, outVal)
			require.False(t, succeed)
			require.Equal(t, "", keyOut)
			require.Equal(t, mockedOrigin, origin)
		},
	)

	// ///////////////////////////////////////////////////////////

	t.Run(
		"slice success and set", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs []int
			dstValue := reflect.ValueOf(rs)

			key := "k4_slice_int"
			origin, keyOut, succeed, outVal, err := b.Bind(key, true, dstValue)
			require.NoError(t, err)
			require.True(t, succeed)
			require.Equal(t, "extK4", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := outVal.Interface()
			require.IsType(t, rs, di)
			require.Equal(t, []int{1, 2, 3, 4, 5}, di.([]int))
			require.True(t, b.entries[key].bounded)
			require.True(t, b.entries[key].used)
		},
	)

	t.Run(
		"slice success and don't set", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs []int
			dstValue := reflect.ValueOf(rs)

			key := "k4_slice_int"
			origin, keyOut, succeed, outVal, err := b.Bind(key, false, dstValue)
			require.NoError(t, err)
			require.False(t, succeed)
			require.Equal(t, "extK4", keyOut)
			require.Equal(t, mockedOrigin, origin)

			require.Equal(t, reflect.Value{}, outVal)

			require.True(t, b.entries[key].bounded)
			require.False(t, b.entries[key].used)
		},
	)

	t.Run(
		"slice parse error", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs []int
			dstValue := reflect.ValueOf(rs)

			key := "k2_string"
			origin, keyOut, succeed, outVal, err := b.Bind(key, false, dstValue)
			require.ErrorIs(t, err, ErrParse)
			require.Equal(t, reflect.Value{}, outVal)
			require.Equal(t, mockedOrigin, origin)
			require.False(t, succeed)
			require.Equal(t, "extK2", keyOut)
			require.ErrorContains(t, err, mockedOriginName)
			require.ErrorContains(t, err, "extK2")
			require.True(t, b.entries[key].bounded)
			require.False(t, b.entries[key].used)
		},
	)

	t.Run(
		"panic when binding invalid type", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs int
			dstValue := reflect.ValueOf(rs)

			key := "k1_int"
			require.Panics(
				t, func() {
					_, _, _, _, _ = b.Bind(key, false, dstValue)
				},
			)
		},
	)
}

type errT struct {
	err    error
	extKey string
}

func convertToErrors(origin string, el []errT) (errs []error) {
	for _, err := range el {
		errs = append(
			errs,
			fmt.Errorf(
				"%s/%s: %w",
				origin,
				err.extKey,
				err.err,
			),
		)
	}

	return
}

func TestBinder_GetErrors(t *testing.T) {
	mockedOriginName := "mocked"
	mockedOrigin := dsco.Origin(mockedOriginName)

	type fields struct {
		internalOpts internalOpts
		entries      entries
		provider     EntriesProvider
	}

	tests := []struct {
		name     string
		fields   fields
		wantErrs []error
	}{
		{
			name: "no errors",
			fields: fields{
				entries: entries{
					"k2a": &entry{
						Entry: Entry{
							ExternalKey: "k2aExt",
						},
						bounded: true,
						used:    true,
					},
					"k3a": &entry{
						Entry: Entry{
							ExternalKey: "k3aExt",
						},
						bounded: true,
						used:    true,
					},
					"k1a": &entry{
						Entry: Entry{
							ExternalKey: "k1aExt",
						},
						bounded: true,
						used:    true,
					},
					"k2b": &entry{
						Entry: Entry{
							ExternalKey: "k2bExt",
						},
						bounded: true,
						used:    true,
					},
					"k3b": &entry{
						Entry: Entry{
							ExternalKey: "k3bExt",
						},
						bounded: true,
						used:    true,
					},
					"k1b": &entry{
						Entry: Entry{
							ExternalKey: "k1bExt",
						},
						bounded: true,
						used:    true,
					},
				},
			},
			wantErrs: nil,
		},
		{
			name: "catch some errors",
			fields: fields{
				entries: entries{
					"k2a": &entry{
						Entry: Entry{
							ExternalKey: "k2aExt",
						},
						bounded: false,
						used:    false,
					},
					"k3a": &entry{
						Entry: Entry{
							ExternalKey: "k3aExt",
						},
						bounded: true,
						used:    false,
					},
					"k1a": &entry{
						Entry: Entry{
							ExternalKey: "k1aExt",
						},
						bounded: true,
						used:    true,
					},
					"k2b": &entry{
						Entry: Entry{
							ExternalKey: "k2bExt",
						},
						bounded: false,
						used:    false,
					},
					"k3b": &entry{
						Entry: Entry{
							ExternalKey: "k3bExt",
						},
						bounded: true,
						used:    false,
					},
					"k1b": &entry{
						Entry: Entry{
							ExternalKey: "k1bExt",
						},
						bounded: true,
						used:    true,
					},
				},
			},
			wantErrs: convertToErrors(
				mockedOriginName,
				[]errT{
					{
						err:    ErrUnboundKey,
						extKey: "k2aExt",
					},
					{
						err:    ErrOverriddenKey,
						extKey: "k3aExt",
					},
					{
						err:    ErrUnboundKey,
						extKey: "k2bExt",
					},
					{
						err:    ErrOverriddenKey,
						extKey: "k3bExt",
					},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				pr := NewMockStrEntriesProvider(t)
				pr.On("GetOrigin").Return(mockedOrigin).Once()

				s := &Binder{
					entries:  tt.fields.entries,
					provider: pr,
				}

				gotErrs := s.GetPostProcessErrors()
				require.ElementsMatch(t, gotErrs, tt.wantErrs)
			},
		)
	}
}
