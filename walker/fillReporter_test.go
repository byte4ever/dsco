package walker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var errMocked1 = errors.New("mocked error 1")
var errMocked2 = errors.New("mocked error 2")

func TestFillReporterImpl_Failed(t *testing.T) {
	t.Parallel()

	t.Run(
		"nil case",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{}

			require.False(t, r.Failed())
		},
	)

	t.Run(
		"empty case",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors: []error{},
			}

			require.False(t, r.Failed())
		},
	)

	t.Run(
		"some errors",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors: []error{nil},
			}

			require.True(t, r.Failed())
		},
	)
}

func TestFillReporterImpl_ReportError(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{}

			r.ReportError(errMocked1)

			require.Equal(
				t,
				FillErrors{errMocked1},
				r.errors,
			)
		},
	)

	t.Run(
		"accum error",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors: FillErrors{errMocked1},
			}

			r.ReportError(errMocked2)

			require.Equal(
				t,
				FillErrors{errMocked1, errMocked2},
				r.errors,
			)
		},
	)
}

func TestFillReporterImpl_ReportOverride(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state no override",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				locations: map[uint]string{},
			}

			r.ReportOverride(11, "some-location")

			require.Empty(
				t,
				r.errors,
			)
		},
	)

	t.Run(
		"initial state with override",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				locations: map[uint]string{11: "prev-location"},
			}

			r.ReportOverride(11, "some-location")

			require.Len(t, r.errors, 1)
			err := r.errors[0]
			require.ErrorIs(t, err, ErrOverriddenKey)
			require.ErrorContains(t, err, "prev-location")
			require.ErrorContains(t, err, "some-location")
		},
	)

	t.Run(
		"accum state no override",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors:    FillErrors{errMocked1},
				locations: map[uint]string{},
			}

			r.ReportOverride(11, "some-location")

			require.Equal(
				t,
				FillErrors{errMocked1},
				r.errors,
			)
		},
	)

	t.Run(
		"accum state with override",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors:    FillErrors{errMocked1},
				locations: map[uint]string{11: "prev-location"},
			}

			r.ReportOverride(11, "some-location")

			require.Len(t, r.errors, 2)
			require.Equal(t, errMocked1, r.errors[0])
			err := r.errors[1]
			require.ErrorIs(t, err, ErrOverriddenKey)
			require.ErrorContains(t, err, "prev-location")
			require.ErrorContains(t, err, "some-location")
		},
	)

}

func TestFillReporterImpl_Result(t *testing.T) {
	t.Parallel()

	t.Run(
		"some errors",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				errors: FillErrors{errMocked1},
			}

			report, err := r.Result()

			require.Nil(
				t,
				report,
			)

			var e FillErrors
			require.ErrorAs(t, err, &e)
			require.Equal(t, e, FillErrors{errMocked1})
		},
	)

	t.Run(
		"no errors",
		func(t *testing.T) {
			t.Parallel()

			expectedReport := FillReport{
				FillReportEntry{
					path:     "some-path",
					location: "some-location",
				},
			}

			r := &FillReporterImpl{
				report: expectedReport,
			}

			gotReport, err := r.Result()

			require.NoError(
				t,
				err,
			)

			require.Equal(t, expectedReport, gotReport)
		},
	)
}

func TestNewFillReporterImpl(t *testing.T) {
	t.Parallel()

	r := NewFillReporterImpl()

	require.NotNil(t, r)
	require.NotNil(t, r.locations)
	require.Nil(t, r.errors)
	require.Nil(t, r.report)
}

func TestFillReporterImpl_ReportUse(t *testing.T) {
	t.Parallel()

	t.Run(
		"initial state",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				locations: map[uint]string{},
			}

			r.ReportUse(111, "some-path", "some-location")

			require.Equal(
				t,
				map[uint]string{111: "some-location"},
				r.locations,
			)

			require.Len(t, r.report, 1)
			require.Equal(
				t, r.report[0], FillReportEntry{
					path:     "some-path",
					location: "some-location",
				},
			)
		},
	)

	t.Run(
		"accum initial state",
		func(t *testing.T) {
			t.Parallel()

			r := &FillReporterImpl{
				report: FillReport{
					{
						path:     "some-path0",
						location: "some-loc0",
					},
				},
				locations: map[uint]string{
					120: "some-loc0",
				},
			}

			r.ReportUse(111, "some-path", "some-location")

			require.Equal(
				t,
				map[uint]string{
					111: "some-location",
					120: "some-loc0",
				},
				r.locations,
			)

			require.Len(t, r.report, 2)
			require.Equal(
				t, r.report[1], FillReportEntry{
					path:     "some-path",
					location: "some-location",
				},
			)
		},
	)

}

func TestFillReporterImpl_ReportUnused(t *testing.T) {
	t.Parallel()

	r := &FillReporterImpl{}
	r.ReportUnused("some-path")

	require.Len(t, r.errors, 1)
	require.ErrorIs(t, r.errors[0], ErrUninitializedKey)
}

func TestFillErrors_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		f    FillErrors
		want string
	}{
		{
			name: "nil",
			f:    nil,
			want: "",
		},
		{
			name: "one error",
			f:    FillErrors{errMocked1},
			want: `mocked error 1
`,
		},
		{
			name: "some errors",
			f:    FillErrors{errMocked1, errMocked2},
			want: `mocked error 1
mocked error 2
`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				if got := tt.f.Error(); got != tt.want {
					t.Errorf("Error() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
