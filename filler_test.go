package dsco

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFiller_errReport2(t *testing.T) {
	t.Parallel()
	report := newMockReportIface(t)
	report.On("perEntryReport").Once().Return([]error{err1, err2, err3})

	layers := newMockLayersIFace(t)
	layers.On("getPostProcessErrors").Once().Return([]error{err4, err5, err6})

	b := &Filler{
		layers: layers,
		report: report,
	}

	errs := b.processReport()
	require.Equal(t, []error{err1, err2, err3, err4, err5, err6}, errs)
}

func TestNewFiller(t *testing.T) {
	t.Parallel()
	t.Run(
		"no layers provided nil case", func(t *testing.T) {
			t.Parallel()
			b, err := NewFiller()
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"no layers provided empty case", func(t *testing.T) {
			t.Parallel()
			b, err := NewFiller([]Binder{}...)
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()
			b1 := NewMockBinder(t)
			b2 := NewMockBinder(t)
			b3 := NewMockBinder(t)
			l := []Binder{b1, b2, b3}
			b, err := NewFiller(l...)
			require.NotNil(t, b)
			require.NoError(t, err)
			require.Equal(t, layers(l), b.layers)
		},
	)
}

func exactMatchString(s string) interface{} {
	return mock.MatchedBy(
		func(v string) bool {
			return v == s
		},
	)
}

func TestFiller_Fill(t *testing.T) {
	t.Parallel()
	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()
			layers := newMockLayersIFace(t)
			val1 := R("stringValue")
			re1 := ReportEntry{
				Value:       reflect.ValueOf(val1),
				Key:         "sub-key1",
				ExternalKey: "sub-key1",
				Idx:         0,
				Errors:      nil,
			}

			val2 := R(uint8(34))
			re2 := ReportEntry{
				Value:       reflect.ValueOf(val2),
				Key:         "sub-key2",
				ExternalKey: "sub-key2",
				Idx:         0,
				Errors:      nil,
			}

			val3 := R(34)
			re3 := ReportEntry{
				Value:       reflect.ValueOf(val3),
				Key:         "key1",
				ExternalKey: "key1",
				Idx:         0,
				Errors:      nil,
			}

			val4 := R(123.321)
			re4 := ReportEntry{
				Value:       reflect.ValueOf(val4),
				Key:         "key2",
				ExternalKey: "key2",
				Idx:         0,
				Errors:      nil,
			}

			n := time.Now().UTC()
			val5 := R(n)
			re5 := ReportEntry{
				Value:       reflect.ValueOf(val5),
				Key:         "key3",
				ExternalKey: "key3",
				Idx:         0,
				Errors:      nil,
			}

			layers.
				On(
					"bind",
					exactMatchString("sub-key1"),
					mock.MatchedBy(
						func(v reflect.Value) bool {
							return true
						},
					),
				).Return(
				re1,
			).Once()

			layers.
				On(
					"bind",
					exactMatchString("sub-key2"),
					mock.MatchedBy(
						func(v reflect.Value) bool {
							return true
						},
					),
				).Return(
				re2,
			).Once()

			layers.
				On(
					"bind",
					exactMatchString("key1"),
					mock.MatchedBy(
						func(v reflect.Value) bool {
							return true
						},
					),
				).Return(
				re3,
			).Once()

			layers.
				On(
					"bind",
					exactMatchString("key2"),
					mock.MatchedBy(
						func(v reflect.Value) bool {
							return true
						},
					),
				).Return(
				re4,
			).Once()

			layers.
				On(
					"bind",
					exactMatchString("key3"),
					mock.MatchedBy(
						func(v reflect.Value) bool {
							return true
						},
					),
				).Return(
				re5,
			).Once()

			report := newMockReportIface(t)

			collectedKeys := make(map[string]int)
			var collectedKeysOrder []string

			report.
				On(
					"addEntry",
					mock.MatchedBy(
						func(v ReportEntry) bool {
							collectedKeys[v.Key]++
							collectedKeysOrder =
								append(collectedKeysOrder, v.Key)

							return true
						},
					),
				).Return().Times(5)

			report.
				On("perEntryReport").
				Return(nil).
				Once()
			layers.
				On("getPostProcessErrors").
				Return(nil).
				Once()

			f := &Filler{
				layers: layers,
				report: report,
			}

			var dst OkRoot
			errs := f.Fill(&dst)
			require.Len(t, errs, 0)
			require.Len(t, collectedKeysOrder, 5)
			require.Len(t, collectedKeys, 5)
			require.Equal(
				t,
				collectedKeysOrder,
				[]string{
					"sub-key1",
					"sub-key2",
					"key1",
					"key2",
					"key3",
				},
			)
			require.Equal(t, val1, dst.Sub.Key1)
			require.Equal(t, val2, dst.Sub.Key2)
			require.Equal(t, val3, dst.Key1)
			require.Equal(t, val4, dst.Key2)
			require.Equal(t, val5, dst.Key3)
		},
	)

	t.Run(
		"not found", func(t *testing.T) {
			t.Parallel()
			layers := newMockLayersIFace(t)
			val1 := R("stringValue")
			re1 := ReportEntry{
				Value:       reflect.ValueOf(val1),
				Key:         "sub-key1",
				ExternalKey: "sub-key1",
				Idx:         -1,
				Errors:      nil,
			}

			val2 := R(uint8(34))
			re2 := ReportEntry{
				Value:       reflect.ValueOf(val2),
				Key:         "sub-key2",
				ExternalKey: "sub-key2",
				Idx:         0,
				Errors:      nil,
			}

			val3 := R(34)
			re3 := ReportEntry{
				Value:       reflect.ValueOf(val3),
				Key:         "key1",
				ExternalKey: "key1",
				Idx:         -1,
				Errors:      nil,
			}

			val4 := R(123.321)
			re4 := ReportEntry{
				Value:       reflect.ValueOf(val4),
				Key:         "key2",
				ExternalKey: "key2",
				Idx:         0,
				Errors:      nil,
			}

			n := time.Now().UTC()
			val5 := R(n)
			re5 := ReportEntry{
				Value:       reflect.ValueOf(val5),
				Key:         "key3",
				ExternalKey: "key3",
				Idx:         -1,
				Errors:      nil,
			}

			layers.On(
				"bind",
				exactMatchString("sub-key1"),
				mock.MatchedBy(
					func(v reflect.Value) bool {
						return true
					},
				),
			).Return(
				re1,
			).Once()

			layers.On(
				"bind",
				exactMatchString("sub-key2"),
				mock.MatchedBy(
					func(v reflect.Value) bool {
						return true
					},
				),
			).Return(
				re2,
			).Once()

			layers.On(
				"bind",
				exactMatchString("key1"),
				mock.MatchedBy(
					func(v reflect.Value) bool {
						return true
					},
				),
			).Return(
				re3,
			).Once()

			layers.On(
				"bind",
				exactMatchString("key2"),
				mock.MatchedBy(
					func(v reflect.Value) bool {
						return true
					},
				),
			).Return(
				re4,
			).Once()

			layers.On(
				"bind",
				exactMatchString("key3"),
				mock.MatchedBy(
					func(v reflect.Value) bool {
						return true
					},
				),
			).Return(
				re5,
			).Once()

			report := newMockReportIface(t)

			collectedKeys := make(map[string]int)
			var collectedKeysOrder []string

			report.
				On(
					"addEntry",
					mock.MatchedBy(
						func(v ReportEntry) bool {
							if !v.isFound() {
								collectedKeys[v.Key]++
								collectedKeysOrder =
									append(collectedKeysOrder, v.Key)
							}

							return true
						},
					),
				).Return().Times(5)

			report.
				On("perEntryReport").
				Return(nil).
				Once()
			layers.
				On("getPostProcessErrors").
				Return(nil).
				Once()

			f := &Filler{
				layers: layers,
				report: report,
			}

			var dst OkRoot
			errs := f.Fill(&dst)
			require.Len(t, errs, 0)
			require.Len(t, collectedKeysOrder, 3)
			require.Len(t, collectedKeys, 3)
			require.Equal(
				t,
				collectedKeysOrder,
				[]string{"sub-key1", "key1", "key3"},
			)
			require.Nil(t, dst.Sub.Key1)
			require.Equal(t, val2, dst.Sub.Key2)
			require.Nil(t, dst.Key1)
			require.Equal(t, val4, dst.Key2)
			require.Nil(t, dst.Key3)
		},
	)

	t.Run(
		"check struct failure", func(t *testing.T) {
			t.Parallel()
			f := &Filler{}
			errs := f.Fill(&T1Root{})
			require.Len(t, errs, 1)
			require.Error(t, errs[0])
			require.ErrorIs(t, errs[0], ErrRecursiveStruct)

		},
	)
}

func Test_formatIndexSequence(t *testing.T) {
	t.Parallel()

	type args struct {
		indexes []int
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "single idx",
			args: args{
				indexes: []int{123},
			},
			want: "#123",
		},
		{
			name: "2 indexes",
			args: args{
				indexes: []int{123, 4000},
			},
			want: "#123 and #4000",
		},
		{
			name: "3 indexes",
			args: args{
				indexes: []int{123, 4000, 233},
			},
			want: "#123, #4000 and #233",
		},
		{
			name: "many indexes",
			args: args{
				indexes: []int{1, 2, 3, 4},
			},
			want: "#1, #2, #3 and #4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				got := formatIndexSequence(tt.args.indexes)
				require.Equal(
					t,
					tt.want,
					got,
					"formatIndexSequence() = %v, want %v",
					got,
					tt.want,
				)
			},
		)
	}
}

func Test_formatIndexSequence_panics(t *testing.T) {
	t.Parallel()
	require.Panics(
		t, func() {
			formatIndexSequence(nil)
		},
	)
	require.Panics(
		t, func() {
			formatIndexSequence([]int{})
		},
	)
}

func TestNewFiller_checkLayers(t *testing.T) {
	t.Parallel()

	type args struct {
		l []Binder
	}

	tests := []struct {
		name              string
		args              args
		want              *Filler
		wantErr           error
		expectErrContains []string
	}{
		{
			name: "nil case",
			args: args{
				l: nil,
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"no layers"},
		},
		{
			name: "empty case",
			args: args{
				l: []Binder{},
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"no layers"},
		},
		{
			name: "single nil binder",
			args: args{
				l: []Binder{nil},
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"#0", "is", "nil"},
		},
		{
			name: "second binder is nil",
			args: args{
				l: []Binder{NewMockBinder(t), nil},
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"#1", "is", "nil"},
		},
		{
			name: "second binder is nil (middle)",
			args: args{
				l: []Binder{
					NewMockBinder(t),
					nil,
					NewMockBinder(t),
				},
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"#1", "is", "nil"},
		},
		{
			name: "last binder is nil",
			args: args{
				l: []Binder{
					NewMockBinder(t),
					NewMockBinder(t),
					NewMockBinder(t),
					nil,
				},
			},
			want:              nil,
			wantErr:           ErrInvalidLayers,
			expectErrContains: []string{"#3", "is", "nil"},
		},
		{
			name: "multiple binder are nil",
			args: args{
				l: []Binder{
					NewMockBinder(t),
					nil, // #1
					NewMockBinder(t),
					NewMockBinder(t),
					nil, // #4
					nil, // #5
					NewMockBinder(t),
					NewMockBinder(t),
					NewMockBinder(t),
					nil, // #9
					nil, // #10
					nil, // #11
				},
			},
			want:    nil,
			wantErr: ErrInvalidLayers,
			expectErrContains: []string{
				"#1",
				"#4",
				"#5",
				"#9",
				"#10",
				"#11",
				"are",
				"nil",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := NewFiller(tt.args.l...)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
					require.Nil(t, got)
					for _, s := range tt.expectErrContains {
						require.ErrorContains(t, err, s)
					}
				}

				require.Equal(
					t,
					tt.want,
					got,
					"NewFiller() got = %v, want %v",
					got, tt.want,
				)
			},
		)
	}
}
