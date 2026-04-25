//nolint:forbidigo // it's an example

// Package main is a preflight check you can run in CI or on container
// startup to list every config key that has no default. If any are
// missing, the program exits with code 2 so the orchestrator can fail
// the deploy before the service even tries to start.
//
// Run with: go run ./examples/inventory/preflight
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

// Config is what the real service would read. Defaults cover the safe
// stuff; anything sensitive (db credentials, API keys) is intentionally
// left unset so an operator must wire it up.
type Config struct {
	Database *DatabaseConfig `yaml:"database"`
	APIKey   *string         `yaml:"api_key"`
	Port     *int            `yaml:"port"`
	Verbose  *bool           `yaml:"verbose"`
}

// DatabaseConfig has no defaults — the operator picks every value.
type DatabaseConfig struct {
	Host     *string `yaml:"host"`
	Name     *string `yaml:"name"`
	User     *string `yaml:"user"`
	Password *string `yaml:"password"`
}

func main() {
	defaults := &Config{
		Port:    dsco.R(8080),
		Verbose: dsco.R(false),
	}

	var cfg *Config

	report, err := inventory.Compute(
		&cfg,
		dsco.WithStructLayer(defaults, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
	)
	if err != nil {
		log.Fatal(err)
	}

	missing := requiredKeys(report)
	if len(missing) == 0 {
		fmt.Println("preflight: all required keys have defaults")
		return
	}

	fmt.Fprintln(os.Stderr, "preflight: required keys with no default:")
	for _, key := range missing {
		fmt.Fprintln(os.Stderr, "  - "+key)
	}
	os.Exit(2)
}

// requiredKeys returns the canonical key for every field that has no
// baked-in default. These are the keys the operator must set for the
// service to start.
func requiredKeys(report *inventory.Report) []string {
	var keys []string
	for _, field := range report.Fields {
		if field.Satisfied != nil {
			continue
		}
		if field.Key == nil {
			keys = append(keys, field.Path)
			continue
		}
		keys = append(keys, field.Key.Key)
	}
	return keys
}
