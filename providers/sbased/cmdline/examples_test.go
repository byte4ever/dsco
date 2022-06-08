package cmdline

import (
	"fmt"
	"sort"
)

func Example_newEntriesProvider_1() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"some-thing"}

	_, err := NewEntriesProvider(params)

	fmt.Println("when processing invalid params got error:")
	fmt.Println(err)
	// Output:
	// when processing invalid params got error:
	// arg #0 - (some-thing): options param not in --xxx=val format
}

func Example_newEntriesProvider_2() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"--arg1=1", "--arg2=1", "--arg1=asdasd"}

	_, err := NewEntriesProvider(params)

	fmt.Println("when processing duplicated params got error:")
	fmt.Println(err)
	// Output:
	// when processing duplicated params got error:
	// --arg1: duplicate param
}

func Example_newEntriesProvider_3() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"--arg3=1", "--arg1=10", "--arg2=asdasd"}

	provider, _ := NewEntriesProvider(params)

	fmt.Println("when processing params got those entries:")

	entries := provider.GetEntries()
	sortedKeys := make([]string, 0, len(entries))

	for key := range entries {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for i, key := range sortedKeys {
		entry := entries[key]
		fmt.Printf(
			"%2d - key=%q externalKey=%q value=%q\n",
			i,
			key,
			entry.ExternalKey,
			entry.Value,
		)
	}
	// Output:
	// when processing params got those entries:
	//  0 - key="arg1" externalKey="--arg1" value="10"
	//  1 - key="arg2" externalKey="--arg2" value="asdasd"
	//  2 - key="arg3" externalKey="--arg3" value="1"
}
