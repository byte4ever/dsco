package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	gc "github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/providers/sbased"
	"github.com/byte4ever/dsco/providers/sbased/cmdline"
	"github.com/byte4ever/dsco/providers/sbased/env"
	"github.com/byte4ever/dsco/providers/strukt"
)

type Zonk struct {
	FirstName    *string
	LastName     *string
	TrainingTime *time.Duration
	T            *time.Time
	B            *bool
}

type Root struct {
	A    *int
	B    *float64
	T    *time.Time
	Z    *Zonk
	L    []string
	NaNa *int
}

func Test(t *testing.T) {
	c1 := &Root{
		L: []string{"A", "B", "C"},
		// A: gc.R(12),
		B: gc.R(123123.123),
		T: gc.R(
			time.Date(
				2014,
				1,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
		),
		Z: &Zonk{
			LastName:     gc.R("ORI2"),
			B:            gc.R(false),
			TrainingTime: gc.R(time.Second),
			T: gc.R(
				time.Date(
					1970,
					2,
					26,
					7,
					0,
					0,
					0,
					time.UTC,
				),
			),
		},
	}

	c2 := &Root{
		L: []string{"LA", "LB"},
		B: gc.R(-20.0),
		Z: &Zonk{
			// FirstName: gc.R("Laurent"),
		},
	}

	c4 := &Root{
		A: gc.R(111),
		Z: &Zonk{
			FirstName: gc.R("NoWay"),
		},
	}

	provideC1, err := strukt.NewBinder(c1)
	require.NoError(t, err)

	provideC2, err := strukt.NewBinder(c2)
	require.NoError(t, err)

	provideC4, err := strukt.NewBinder(c4)
	require.NoError(t, err)

	require.NoError(t, os.Setenv("SRV-VERBOSITY", `yes`))
	require.NoError(t, os.Setenv("SRV-L", `[q, w, e]`))
	require.NoError(t, os.Setenv("SRV-NA_NA", `111`))

	provideC3p, errs := env.NewEntriesProvider("SRV")
	require.Nil(t, errs)

	provideC3, err := sbased.NewBinder(
		provideC3p,
		sbased.WithAliases(
			map[string]string{
				"verbosity": "z-b",
			},
		),
	)
	require.NoError(t, err)

	provideC5p, err := cmdline.NewEntriesProvider(
		[]string{
			// `--shitty=123123`,
			`--last_name=DOULOS`,
			`--z-training_time=23s`,
			// `--a=1000`,
		},
	)
	require.NoError(t, err)

	provideC5, err := sbased.NewBinder(
		provideC5p,
		sbased.WithAliases(
			map[string]string{
				"last_name": "z-last_name",
				"shitty":    "should-not-scan",
			},
		),
	)
	require.NoError(t, err)

	cf, err := gc.NewFiller(
		provideC4,
		provideC5,
		provideC3,
		provideC2,
		provideC1,
	)
	require.NoError(t, err)

	var cc Root
	errs = cf.Fill(&cc)
	require.Empty(t, errs)

	ll, err := yaml.Marshal(cc)
	require.NoError(t, err)
	t.Log("---------------------------------------------")
	t.Log(string(ll))
	t.Log("---------------------------------------------")
}
