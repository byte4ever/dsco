package cmdline

import (
	"testing"
)

func FuzzNewEntriesProvider(f *testing.F) {
	f.Add("--arg1=asdasd", "--arg1=asdasd1")
	f.Add("--arg1_-asd=asdasd", "--arg2=asdasd1")
	f.Add("--arg1-asd", "--arg2=asdasd1")
	f.Fuzz(
		func(t *testing.T, param1 string, param2 string) {
			p, err := NewEntriesProvider([]string{param1, param1})
			if p != nil && err != nil {
				t.Errorf("%v %v", p, err)
			}
		},
	)
}
