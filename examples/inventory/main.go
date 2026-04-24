//nolint:forbidigo // it's an example

// Package main demonstrates the dsco inventory sub-package.
//
// Run with: go run ./examples/inventory.
package main

import (
	"log"
	"os"

	"github.com/byte4ever/dsco"
	"github.com/byte4ever/dsco/inventory"
)

// Config is a minimal example configuration showing several layer types.
type Config struct {
	Host    *string `yaml:"host"`
	Port    *int    `yaml:"port"`
	Verbose *bool   `yaml:"verbose"`
}

func main() {
	defaults := &Config{
		Port:    dsco.R(8080),
		Verbose: dsco.R(false),
	}

	var c *Config

	report, err := inventory.Compute(
		&c,
		dsco.WithStructLayer(defaults, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = report.WriteText(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
