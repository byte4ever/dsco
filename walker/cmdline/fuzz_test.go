package cmdline

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func FuzzNewEntriesProvider(f *testing.F) {
	testsParams := [][]string{
		{},
		{"--arg1=asdasd", "--arg1=asdasd1", "--arg3=asdasd", "--arg4=asdasd1"},
		{"--arg1=asdasd"},
		{"--arg1_-asd=asdasd", "--arg2=asdasd1"},
		{"--arg1_-asd=asdasd", "--arg2=asdasd1"},
	}

	for _, param := range testsParams {
		data, _ := yaml.Marshal(param)
		f.Add(data)
	}

	f.Fuzz(
		func(t *testing.T, data []byte) {
			var vr []string

			if err := yaml.Unmarshal(data, &vr); err != nil {
				t.Skip()
			}
			t.Log(vr)
			p, err := NewEntriesProvider(vr)
			if p != nil && err != nil {
				t.Errorf("%v %v", p, err)
			}
		},
	)
}
