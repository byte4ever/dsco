package env

import (
	"fmt"
	"testing"
)

func FuzzNewEntriesProvider(f *testing.F) {
	f.Add("PREFIX", "PREFIX_A_B-C", "value")
	f.Add("PREFIX", "PREFIX_A", "value")
	f.Add("PREFIX", "PREFIX_A_B-DC", "value")
	f.Add("PREFIX", "P_A_B-DC", "value")
	f.Add("PREFIX", "PREFIX_A_B-DC", "value=asdasd")
	f.Add("PREFIX", "PREFIX_A_B-_DC", "value=asdasd")
	f.Add("PREFIX", "PREFIX_A_1B-DC", "value=asdasd")
	f.Add("PREFIX", "PREFIX_A_B-1DC", "value=asdasd")
	f.Fuzz(
		func(t *testing.T, prefix string, key string, value string) {
			envKeyVar := fmt.Sprintf("%s=%s", key, value)
			p, errs := newProvider(prefix, []string{envKeyVar})
			if p != nil && errs != nil {
				t.Errorf("%v %v", p, errs)
			}
		},
	)
}
