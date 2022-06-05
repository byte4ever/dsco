package env

import (
	"fmt"
	"os"
	"sort"
)

func ExampleNewEntriesProvider() {
	_, err := NewEntriesProvider("asd")
	fmt.Printf("Got error: %v\n", err)

	// creates some env var for testing purpose
	err = os.Setenv("API-SOME-KEY1", "value1")
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

	// Note the provider will only grab the environment variables that starts with
	// the given prefix.
	p, err := NewEntriesProvider("API")
	if err != nil {
		panic(err)
	}

	fmt.Println("\nHere are the entries")

	entries := p.GetEntries()

	keys := make([]string, 0, len(entries))

	for s := range entries {
		keys = append(keys, s)
	}

	sort.Strings(keys)

	for _, key := range keys {
		entry := entries[key]
		fmt.Printf(
			"key:<%s>  externalKey:<%s> value:<%s>\n",
			key,
			entry.ExternalKey,
			entry.Value,
		)
	}

	// Output:
	// Got error: "asd" : invalid prefix
	//
	// Here are the entries
	// key:<some-key1>  externalKey:<API-SOME-KEY1> value:<value1>
	// key:<some-key2>  externalKey:<API-SOME-KEY2> value:<value2>
	// key:<some-key3>  externalKey:<API-SOME-KEY3> value:<value3>
}
