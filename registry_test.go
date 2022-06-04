package dsco

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var err1 = errors.New("mocked error 1")
var err2 = errors.New("mocked error 2")
var err3 = errors.New("mocked error 3")
var err4 = errors.New("mocked error 4")
var err5 = errors.New("mocked error 5")
var err6 = errors.New("mocked error 6")
var err7 = errors.New("mocked error 7")
var err8 = errors.New("mocked error 8")
var err9 = errors.New("mocked error 9")

func buildBinder(t *testing.T, errs []error) Binder {
	n := NewMockBinder(t)
	n.On("GetPostProcessErrors").
		Once().Return(errs)

	return n
}

func getLayers(t *testing.T, mr ...[]error) (bs []Binder) {
	t.Helper()

	for _, errs := range mr {
		bs = append(bs, buildBinder(t, errs))
	}

	return
}

func TestFiller_errReport2(t *testing.T) {
	layersErrors := [][]error{{err1}, {}, {err2, err3}, {}, {err4, err5, err6}}

	var layers Layers

	for _, layersError := range layersErrors {
		l := getLayers(t, layersError)
		layers = append(layers, l...)
	}

	b := &Filler{
		layers: layers,
		m: Report{
			ReportEntry{
				Key:         "invalidKey",
				ExternalKey: "extKey",
				Idx:         -1,
			},
			ReportEntry{
				Errors: []error{err7, err8},
			},
			ReportEntry{
				Errors: []error{},
			},
			ReportEntry{
				Errors: []error{err9},
			},
		},
	}

	errs := b.processReport()
	require.Equal(t, []error{err7, err8, err9, err1, err2, err3, err4, err5, err6}, errs[1:])

	ur := errs[0]
	require.ErrorIs(t, ur, ErrUninitialized)
	require.ErrorContains(t, ur, "invalidKey")
}

func TestNewFiller(t *testing.T) {
	t.Run(
		"no layers provided nil case", func(t *testing.T) {
			b, err := NewFiller(nil)
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"no layers provided empty case", func(t *testing.T) {
			b, err := NewFiller([]Binder{})
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"success", func(t *testing.T) {
			b1 := NewMockBinder(t)
			b2 := NewMockBinder(t)
			b3 := NewMockBinder(t)
			layers := []Binder{b1, b2, b3}
			b, err := NewFiller(layers)
			require.NotNil(t, b)
			require.NoError(t, err)
			require.Equal(t, Layers(layers), b.layers)
		},
	)
}
