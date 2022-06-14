package env

import (
	"fmt"
	"os"
)

func Example_newEntriesProvider_1() {
	// Invalid Prefix detection
	_, err := NewEntriesProvider("asd")

	if err == nil {
		panic("err must be returned")
	}

	fmt.Println("Got error:")
	fmt.Println(err)

	// Output:
	// Got error:
	// "asd" : invalid prefix
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

	provider, err := NewEntriesProvider("API")

	if err != nil {
		panic(err)
	}

	fmt.Println("\nHere are the stringValues")

	stringValues := provider.GetStringValues()

	for i, value := range stringValues {
		fmt.Printf(
			"%2d - key:<%s>  location:<%s> value:<%s>\n",
			i,
			value.Key,
			value.Location,
			value.Value,
		)
	}
	// Output:
	// Here are the stringValues
	//  0 - key:<some-key1>  location:<env[API-SOME-KEY1]> value:<value1>
	//  1 - key:<some-key2>  location:<env[API-SOME-KEY2]> value:<value2>
	//  2 - key:<some-key3>  location:<env[API-SOME-KEY3]> value:<value3>
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

	_, err = NewEntriesProvider("API")
	if err == nil {
		panic("must return an error")
	}

	fmt.Println("Got error:")
	fmt.Println(err)
	// Output:
	// Got error:
	// "API-SOME-_KEY2", "API-SOME-kEY3" and "APISOME-KEY1" are ambiguous
}
