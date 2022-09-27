//nolint:forbidigo // it's an example

package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
)

type RetryConfTmpl struct {
	BackOffFactor *float64 `yaml:"back_off_factor"`
	Retry         *int     `yaml:"retry"`
}

type HttpBasedConfTmpl struct {
	RetryConfTmpl `yaml:"retry"`
	URL           *string `yaml:"url"`
	Verbose       *bool   `yaml:"verbose"`
}

type AuthentServiceConf struct {
	HttpBasedConfTmpl `yaml:"http"`
	AccessToken       *string `yaml:"access_token"`
}

type ClientAPIConf struct {
	HttpBasedConfTmpl `yaml:"http"`
	EnableSecurity    *bool `yaml:"enable_security"`
}

type MainConf struct {
	Authentication *AuthentServiceConf `yaml:"authentication"`
	ClientAPI      *ClientAPIConf      `yaml:"client_api"`
	PingDuration   *time.Duration      `yaml:"ping_duration"`
}

func main() {

	// DSCO will try to fill (and allocate the config struct
	var pp *MainConf
	fillReport, err := dsco.Fill(
		// provide a ref
		&pp,

		// Only one command line can be present
		//
		// struct path will be mapped this way
		// Authentication.AccessToken -> --authentication-access_token
		//
		// You can use aliases see next layer.
		dsco.WithCmdlineLayer(),

		// Matches any env var
		//
		// Previous layer cannot override it.
		//
		// You can add multiple layers with different prefixes that's up to
		// you...
		dsco.WithStrictEnvLayer(
			"SRV",
			dsco.WithAliases(
				map[string]string{
					// can use SRV-TOKEN instead of the long version by
					// defining an alias.
					"token": "authentication-access_token",
				},
			),
		),

		// Matches the given go struct
		//
		// No values here can be overridden by the
		// previous layer. Not even the previous env layer
		dsco.WithStrictStructLayer(
			&MainConf{
				//  let say that authentication is hardcoded
				Authentication: &AuthentServiceConf{
					HttpBasedConfTmpl: HttpBasedConfTmpl{
						RetryConfTmpl: RetryConfTmpl{
							BackOffFactor: dsco.R(1.2),
							Retry:         dsco.R(5),
						},
						URL: dsco.R("http://perfect-authent.com"),
					},
				},
			},
			"immutable", // <- this is the layer id
		),

		// This layer defines values that can be overridden by all previous
		// layers.
		//
		// So it acts as a kind of fallback layer.
		dsco.WithStructLayer(
			&MainConf{
				Authentication: &AuthentServiceConf{
					HttpBasedConfTmpl: HttpBasedConfTmpl{
						// set some default retry
						RetryConfTmpl: RetryConfTmpl{
							BackOffFactor: dsco.R(1.05),
							Retry:         dsco.R(20),
						},
						// verbosity is false by default
						Verbose: dsco.R(false),
					},
				},
				ClientAPI: &ClientAPIConf{
					HttpBasedConfTmpl: HttpBasedConfTmpl{
						// set some default retry
						RetryConfTmpl: RetryConfTmpl{
							BackOffFactor: dsco.R(1.05),
							Retry:         dsco.R(20),
						},
						// verbosity is false by default
						Verbose: dsco.R(false),
					},
				},
				// ping duration is 10s by default
				PingDuration: dsco.R(10 * time.Second),
			},
			"mutable", // <- this is the layer id
		),
	)

	// If structure fill fails because of missing value field then structure
	// is partially filled (i.e pointer is not nil).
	// This is might be useful for debugging purpose.
	if pp != nil {
		fmt.Println("filled structure ____________________")
		s, _ := yaml.Marshal(pp)
		fmt.Println(string(s))
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		// pp is completely filled (i.e all fields are defined).

		fmt.Println("\nfill report for debugging purpose____")
		fmt.Println(fillReport)
		fillReport.Dump(os.Stdout)
	}

}
