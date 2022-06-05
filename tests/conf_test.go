package tests

import (
	"crypto"
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
	"github.com/byte4ever/dsco/utils/hash"
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
	H    *hash.Hash
	T    *time.Time
	Z    *Zonk
	L    []string
	NaNa *int
}

func Test(t *testing.T) {
	c1 := &Root{
		L: []string{"A", "B", "C"},
		// A: gc.V(12),
		B: gc.V(123123.123),
		H: gc.V(hash.Hash(crypto.SHA256)),
		T: gc.V(
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
			LastName:     gc.V("ORI2"),
			B:            gc.V(false),
			TrainingTime: gc.V(time.Second),
			T: gc.V(
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
		B: gc.V(-20.0),
		Z: &Zonk{
			// FirstName: gc.V("Laurent"),
		},
	}

	c4 := &Root{
		A: gc.V(111),
		Z: &Zonk{
			FirstName: gc.V("NoWay"),
		},
	}

	provideC1, err := strukt.Provide(c1)
	require.NoError(t, err)

	provideC2, err := strukt.Provide(c2)
	require.NoError(t, err)

	provideC4, err := strukt.Provide(c4)
	require.NoError(t, err)

	// require.NoError(t, os.Setenv("SRV-Z-FIRST_NAME", `Celina`))
	require.NoError(t, os.Setenv("SRV-VERBOSITY", `yes`))
	require.NoError(t, os.Setenv("SRV-L", `[q, w, e]`))
	// require.NoError(t, os.Setenv("SRV-BITOS", `asd`))
	require.NoError(t, os.Setenv("SRV-NA_NA", `111`))

	provideC3p, err := env.Provide("SRV")
	require.NoError(t, err)

	provideC3, err := sbased.Provide(
		provideC3p,
		sbased.WithAliases(
			map[string]string{
				"verbosity": "z-b",
			},
		),
	)
	require.NoError(t, err)

	provideC5p, err := cmdline.Provide(
		[]string{
			// `--shitty=123123`,
			`--last_name=DOULOS`,
			`--z-training_time=23s`,
			// `--a=1000`,
		},
	)
	require.NoError(t, err)

	provideC5, err := sbased.Provide(
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
	errs := cf.Fill(&cc)
	require.Empty(t, errs)

	ll, err := yaml.Marshal(cc)
	require.NoError(t, err)
	t.Log("---------------------------------------------")
	t.Log(string(ll))
	t.Log("---------------------------------------------")
}
