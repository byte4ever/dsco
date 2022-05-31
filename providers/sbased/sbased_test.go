package sbased

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/goconf"
)

func TestProvide(t *testing.T) {
	t.Run(
		"success", func(t *testing.T) {
			sep := NewMockStrEntriesProvider(t)
			sep.On("GetEntries").Once().
				Return(
					StrEntries{
						"k1": &StrEntry{
							ExternalKey: "extK1",
							Value:       "val1",
						},
						"k2": &StrEntry{
							ExternalKey: "extK2",
							Value:       "val2",
						},
					},
				)

			p, err := Provide(sep)
			require.NoError(t, err)
			require.NotNil(t, p)

			require.Equal(
				t, p, &Binder{
					internalOpts: internalOpts{
						aliases: nil,
					},
					entries: entries{
						"k1": &entry{
							StrEntry: StrEntry{
								ExternalKey: "extK1",
								Value:       "val1",
							},
							bounded: false,
							used:    false,
						},
						"k2": &entry{
							StrEntry: StrEntry{
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
					StrEntries{
						"alias1": &StrEntry{
							ExternalKey: "extK1Aliased",
							Value:       "val1",
						},
						"k2": &StrEntry{
							ExternalKey: "extK2",
							Value:       "val2",
						},
						"alias3": &StrEntry{
							ExternalKey: "extK3Aliased",
							Value:       "val3",
						},
					},
				)

			p, err := Provide(
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
							StrEntry: StrEntry{
								ExternalKey: "extK1Aliased",
								Value:       "val1",
							},
							bounded: false,
							used:    false,
						},
						"k2": &entry{
							StrEntry: StrEntry{
								ExternalKey: "extK2",
								Value:       "val2",
							},
							bounded: false,
							used:    false,
						},
						"k3": &entry{
							StrEntry: StrEntry{
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

			p, err := Provide(
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
				StrEntry: StrEntry{
					ExternalKey: "extK1Aliased",
					Value:       "1234",
				},
				bounded: false,
				used:    false,
			},
			"k2_string": &entry{
				StrEntry: StrEntry{
					ExternalKey: "extK2",
					Value:       "val2",
				},
				bounded: false,
				used:    false,
			},
			"k3_float64": &entry{
				StrEntry: StrEntry{
					ExternalKey: "extK3Aliased",
					Value:       "val3",
				},
				bounded: false,
				used:    false,
			},
			"k4_slice_int": &entry{
				StrEntry: StrEntry{
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
	mockedOrigin := goconf.Origin(mockedOriginName)

	t.Run(
		"success and set", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, err := b.Bind("k1_int", true, dstType, &dstValue)
			require.NoError(t, err)
			require.True(t, succeed)
			require.Equal(t, "extK1Aliased", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := dstValue.Interface()
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, err := b.Bind("k1_int", false, dstType, &dstValue)
			require.NoError(t, err)
			require.False(t, succeed)
			require.Equal(t, "extK1Aliased", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := dstValue.Interface()
			require.IsType(t, rs, di)
			require.Equal(t, (*int)(nil), di.(*int))
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "k2_string"
			origin, keyOut, succeed, err := b.Bind(key, true, dstType, &dstValue)
			require.ErrorIs(t, err, ErrParse)
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "alias1"
			origin, keyOut, succeed, err := b.Bind(key, true, dstType, &dstValue)
			require.ErrorIs(t, err, ErrAliasCollision)
			require.Equal(t, mockedOrigin, origin)
			require.False(t, succeed)
			require.Equal(t, "", keyOut)
			require.ErrorContains(t, err, mockedOriginName)
			require.ErrorContains(t, err, key)
			// require.False(t, b.entries[key].bounded)
			// require.False(t, b.entries[key].used)
		},
	)

	t.Run(
		"key not found", func(t *testing.T) {
			pr := NewMockStrEntriesProvider(t)
			pr.On("GetOrigin").Return(mockedOrigin).Once()
			b := getBinder()
			b.provider = pr

			var rs *int
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			origin, keyOut, succeed, err := b.Bind("not_found", true, dstType, &dstValue)
			require.NoError(t, err)
			require.False(t, succeed)
			require.Equal(t, "", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := dstValue.Interface()
			require.IsType(t, rs, di)
			require.Nil(t, di.(*int))
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "k4_slice_int"
			origin, keyOut, succeed, err := b.Bind(key, true, dstType, &dstValue)
			require.NoError(t, err)
			require.True(t, succeed)
			require.Equal(t, "extK4", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := dstValue.Interface()
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "k4_slice_int"
			origin, keyOut, succeed, err := b.Bind(key, false, dstType, &dstValue)
			require.NoError(t, err)
			require.False(t, succeed)
			require.Equal(t, "extK4", keyOut)
			require.Equal(t, mockedOrigin, origin)

			di := dstValue.Interface()
			require.IsType(t, rs, di)
			require.Nil(t, di.([]int))
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "k2_string"
			origin, keyOut, succeed, err := b.Bind(key, false, dstType, &dstValue)
			require.ErrorIs(t, err, ErrParse)
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
			dstType := reflect.TypeOf(rs)
			dstValue := reflect.ValueOf(rs)

			key := "k1_int"
			require.Panics(
				t, func() {
					_, _, _, _ = b.Bind(key, false, dstType, &dstValue)
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
	mockedOrigin := goconf.Origin(mockedOriginName)

	type fields struct {
		internalOpts internalOpts
		entries      entries
		provider     StrEntriesProvider
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
						StrEntry: StrEntry{
							ExternalKey: "k2aExt",
						},
						bounded: true,
						used:    true,
					},
					"k3a": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k3aExt",
						},
						bounded: true,
						used:    true,
					},
					"k1a": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k1aExt",
						},
						bounded: true,
						used:    true,
					},
					"k2b": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k2bExt",
						},
						bounded: true,
						used:    true,
					},
					"k3b": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k3bExt",
						},
						bounded: true,
						used:    true,
					},
					"k1b": &entry{
						StrEntry: StrEntry{
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
						StrEntry: StrEntry{
							ExternalKey: "k2aExt",
						},
						bounded: false,
						used:    false,
					},
					"k3a": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k3aExt",
						},
						bounded: true,
						used:    false,
					},
					"k1a": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k1aExt",
						},
						bounded: true,
						used:    true,
					},
					"k2b": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k2bExt",
						},
						bounded: false,
						used:    false,
					},
					"k3b": &entry{
						StrEntry: StrEntry{
							ExternalKey: "k3bExt",
						},
						bounded: true,
						used:    false,
					},
					"k1b": &entry{
						StrEntry: StrEntry{
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
