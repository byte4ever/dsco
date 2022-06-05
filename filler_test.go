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

func TestFiller_errReport2(t *testing.T) {
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
	t.Run(
		"no layers provided nil case", func(t *testing.T) {
			b, err := NewFiller()
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"no layers provided empty case", func(t *testing.T) {
			b, err := NewFiller([]Binder{}...)
			require.Nil(t, b)
			require.ErrorIs(t, err, ErrInvalidLayers)
		},
	)

	t.Run(
		"success", func(t *testing.T) {
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
