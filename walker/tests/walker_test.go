package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/walker"
)

type Sub struct {
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
	Z    *Sub
	NaNa *int
	L    []string
}

func Test_lab2(t *testing.T) { //nolint:paralleltest // using setenv
	os.Args = []string{
		"appName",
		"--a=1234",
		"--z-b=yes",
		"--z-first_name=Laura",
	}

	// t.Setenv("TST-A", "123")
	t.Setenv("TST-B", "123.1234")
	// t.Setenv("API-Z-FIRST_NAME", "Laurent")

	var pp *Root
	fillReport, err := walker.Fill(
		&pp,
		walker.WithEnvLayer("API"),
		walker.WithEnvLayer("TST"),
		walker.WithStrictCmdlineLayer(),
		walker.WithStructLayer(
			&Root{
				B: dsco.R(0.0),
				Z: &Sub{
					FirstName: dsco.R("Rose"),
					LastName:  dsco.R("Dupont"),
					B:         dsco.R(false),
				},
			}, "dflt1",
		),
		walker.WithStructLayer(
			&Root{
				A: dsco.R(120),
				B: dsco.R(2333.32),
				T: dsco.R(time.Now().UTC()),
				Z: &Sub{
					FirstName:    dsco.R("Lola"),
					LastName:     dsco.R("MARTIN"),
					TrainingTime: dsco.R(800 * time.Second),
					T:            dsco.R(time.Now().UTC()),
					B:            dsco.R(true),
				},
				NaNa: dsco.R(2331),
				L:    []string{"A", "B", "C"},
			}, "dflt2",
		),
	)

	if err != nil {
		t.Log(err)
	} else {
		bb, err := yaml.Marshal(pp)
		require.NoError(t, err)

		t.Log(string(bb))

		fillReport.Dump(os.Stdout)
	}
}
