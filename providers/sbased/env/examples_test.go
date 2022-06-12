package env

import (
	"fmt"
	"os"
	"sort"
)

func Example_newEntriesProvider_1() {
	// Invalid Prefix detection
	_, errs := NewEntriesProvider("asd")

	if len(errs) == 0 {
		panic("")
	}

	fmt.Println("Got errors:")

	for _, err := range errs {
		fmt.Println("-", err)
	}

	// Output:
	// Got errors:
	// - "asd" : invalid prefix
}

func Example_newEntriesProvider_2() {
	// creates some env var for testing purpose only
	err := os.Setenv("API-SOME-KEY1", "value1")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("API-SOME-KEY2", "value2")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("API-SOME-KEY3", "value3")

	if err != nil {
		panic(err)
	}

	err = os.Setenv("OTHERPREFIX-SOME-KEY3", "value3")
	if err != nil {
		panic(err)
	}

	// N.B the provider will only grab the environment variables that starts
	// with the provided prefix.

	provider, errs := NewEntriesProvider("API")

	if errs != nil {
		panic(err)
	}

	fmt.Println("\nHere are the entries")

	entries := provider.GetEntries()

	sortedKeys := make([]string, 0, len(entries))

	for s := range entries {
		sortedKeys = append(sortedKeys, s)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		entry := entries[key]
		fmt.Printf(
			"key:<%s>  externalKey:<%s> value:<%s>\n",
			key,
			entry.ExternalKey,
			entry.Value,
		)
	}

	// Output:
	// Here are the entries
	// key:<some-key1>  externalKey:<API-SOME-KEY1> value:<value1>
	// key:<some-key2>  externalKey:<API-SOME-KEY2> value:<value2>
	// key:<some-key3>  externalKey:<API-SOME-KEY3> value:<value3>
}

func Example_newEntriesProvider_3() {
	// creates some env var for testing purpose only
	err := os.Setenv("APISOME-KEY1", "value1")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("API-SOME-_KEY2", "value2")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("API-SOME-kEY3", "value3")
	if err != nil {
		panic(err)
	}

	// that key will not be scanned because prefix does not match.
	err = os.Setenv("OTHERPREFIX-SOME-KEY3", "value3")
	if err != nil {
		panic(err)
	}

	// that key will match perfectly, but malformed keys will abort the
	// creation.
	err = os.Setenv("API-SOME-KEY4", "value3")
	if err != nil {
		panic(err)
	}

	p, errs := NewEntriesProvider("API")

	if p != nil {
		panic("")
	}

	fmt.Println("Got errors:")

	for _, err := range errs {
		fmt.Println("-", err)
	}

	// Output:
	// Got errors:
	// - env var API-SOME-_KEY2: invalid key Format
	// - env var API-SOME-kEY3: invalid key Format
	// - env var APISOME-KEY1: invalid key Format
}
