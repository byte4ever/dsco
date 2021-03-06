package plocation

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathLocations_Dump(t *testing.T) {
	t.Parallel()

	t.Run(
		"", func(t *testing.T) {
			t.Parallel()

			ploc := Locations{
				{
					UID:      0,
					Path:     "path0",
					Location: "loc0",
				},
				{
					UID:      1,
					Path:     "path1",
					Location: "loc1",
				},
			}

			v := bytes.NewBufferString("")

			ploc.Dump(v)

			expectedString := `  path   |  Location
  ----   |  --------
  path0  |  loc0
  path1  |  loc1
`

			require.Equal(t, expectedString, v.String())
		},
	)
}

func TestPathLocations_Report(t *testing.T) {
	t.Parallel()

	ploc := Locations{
		{
			UID:      0,
			Path:     "path0",
			Location: "loc0",
		},
	}

	ploc.Report(1, "path1", "loc1")
	ploc.Report(2, "path2", "loc2")

	require.Equal(
		t, Locations{
			{
				UID:      0,
				Path:     "path0",
				Location: "loc0",
			},
			{
				UID:      1,
				Path:     "path1",
				Location: "loc1",
			},
			{
				UID:      2,
				Path:     "path2",
				Location: "loc2",
			},
		}, ploc,
	)
}

func TestPathLocations_ReportOther(t *testing.T) {
	t.Parallel()

	ploc1 := Locations{
		{
			UID:      0,
			Path:     "path0",
			Location: "loc0",
		},
	}

	ploc2 := Locations{
		{
			UID:      1,
			Path:     "path1",
			Location: "loc1",
		},
		{
			UID:      2,
			Path:     "path2",
			Location: "loc2",
		},
	}

	ploc1.Append(ploc2)

	require.Equal(
		t, Locations{
			{
				UID:      0,
				Path:     "path0",
				Location: "loc0",
			},
			{
				UID:      1,
				Path:     "path1",
				Location: "loc1",
			},
			{
				UID:      2,
				Path:     "path2",
				Location: "loc2",
			},
		}, ploc1,
	)
}
