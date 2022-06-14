package cmdline

import (
	"fmt"
)

func Example_newEntriesProvider_1() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"some-thing"}

	_, err := NewEntriesProvider(params)

	fmt.Println("when processing invalid params got error:")
	fmt.Println(err)
	// Output:
	// when processing invalid params got error:
	// error found at position #1: arg some-thing: options param not in --xxx=val format
}

func Example_newEntriesProvider_2() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"--arg1=1", "--arg2=1", "--arg1=asdasd"}

	_, err := NewEntriesProvider(params)

	fmt.Println("when processing duplicated params got error:")
	fmt.Println(err)
	// Output:
	// when processing duplicated params got error:
	// error found at position #3: --arg1 previous found at position #1: duplicate param
}

func Example_newEntriesProvider_3() {
	// defines params line command. os.Arg[1:] is commonly used.
	params := []string{"--arg3=1", "--arg1=10", "--arg2=asdasd"}

	provider, _ := NewEntriesProvider(params)

	fmt.Println("when processing params got those entries:")

	for i, entry := range provider.GetStringValues() {
		fmt.Printf(
			"%2d - location=%q key=%q value=%q\n",
			i,
			entry.Location,
			entry.Key,
			entry.Value,
		)
	}
	// Output:
	// when processing params got those entries:
	//  0 - location="cmdline[--arg3]" key="arg3" value="1"
	//  1 - location="cmdline[--arg1]" key="arg1" value="10"
	//  2 - location="cmdline[--arg2]" key="arg2" value="asdasd"
}
