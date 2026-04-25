//nolint:forbidigo // it's an example

// Package main shows how to dump a dsco inventory as JSON, the format
// you'd feed to jq, ansible, or whatever else fills env files for you.
//
// Run with: go run ./examples/inventory/json | jq
package main

import (
	"log"
	"os"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

// Config holds the settings our service reads at startup. Pointer fields
// let us tell "not set" apart from "set to the zero value".
type Config struct {
	Database *DatabaseConfig `yaml:"database"`
	Port     *int            `yaml:"port"`
	Verbose  *bool           `yaml:"verbose"`
}

// DatabaseConfig is nested so the example shows how paths look for
// fields a couple of levels deep.
type DatabaseConfig struct {
	Host *string `yaml:"host"`
	Name *string `yaml:"name"`
	User *string `yaml:"user"`
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

	if err := report.WriteJSON(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
