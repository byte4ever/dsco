//nolint:forbidigo // it's an example

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/internal/kfile"
)

// RetryConfTmpl is a sample config.
type RetryConfTmpl struct {
	BackOffFactor *float64 `yaml:"back_off_factor"`
	Retry         *int     `yaml:"retry"`
}

// HTTPBasedConfTmpl is a sample config.
type HTTPBasedConfTmpl struct {
	RetryConfTmpl `yaml:"retry"`
	URL           *string `yaml:"url"`
	Verbose       *bool   `yaml:"verbose"`
}

// AuthentServiceConf is a sample config.
type AuthentServiceConf struct {
	HTTPBasedConfTmpl `yaml:",inline"`
	AccessToken       *string `yaml:"access_token"`
}

// ClientAPIConf is a sample config.
type ClientAPIConf struct {
	HTTPBasedConfTmpl `yaml:",inline"`
	EnableSecurity    *bool `yaml:"enable_security"`
}

// MainConf is a sample config.
type MainConf struct {
	Authentication *AuthentServiceConf `yaml:"authentication"`
	ClientAPI      *ClientAPIConf      `yaml:"client_api"`
	PingDuration   *time.Duration      `yaml:"ping_duration"`
	SecretKey1     *string             `yaml:"secret_key1"`
	SecretKey2     *string             `yaml:"secret_key2"`
}

func main() {
	// try to get some secrets from file system
	secretProvider, err := kfile.NewEntriesProvider(
		"examples/deadsimple/secrets")
	if err != nil {
		log.Fatal(err)
	}

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
					HTTPBasedConfTmpl: HTTPBasedConfTmpl{
						RetryConfTmpl: RetryConfTmpl{
							BackOffFactor: dsco.R(1.2),
							Retry:         dsco.R(5),
						},
						URL: dsco.R("is a sample config.is a sample config.http://perfect-authent.com"),
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
					HTTPBasedConfTmpl: HTTPBasedConfTmpl{
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
					HTTPBasedConfTmpl: HTTPBasedConfTmpl{
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
		dsco.WithStringValueProvider(secretProvider),
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
	}

	// pp is completely filled (i.e all fields are defined).
	fmt.Println("\nfill report for debugging purpose____")
	fmt.Println(fillReport)
	fillReport.Dump(os.Stdout)
}
