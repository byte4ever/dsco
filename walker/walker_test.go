package walker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// type ET2 struct {
// 	W  *int `yaml:",omitempty"`
// 	X1 *int `yaml:",omitempty"`
// }
//
// type ET1 struct {
// 	X2  *int `yaml:",omitempty"`
// 	ET2 `yaml:",inline"`
// }
//
// type F1 struct {
// 	X *int `yaml:",omitempty"`
// }
//
// type F2 struct {
// 	Z *int `yaml:",omitempty"`
// }
//
// type MST struct {
// 	X *int
// 	Y *int
// }
//
// type Root struct {
// 	ET1 `yaml:",inline"`
// 	A   *int `yaml:",omitempty"`
// 	B   *int `yaml:",omitempty"`
// 	P   *int `yaml:",omitempty"`
// 	F1  `yaml:",inline"`
// 	F2  `yaml:",inline"`
// 	M   *MST
// }
//

type Zonk struct {
	FirstName    *string
	LastName     *string
	TrainingTime *time.Duration
	T            *time.Time
	B            *bool
}

type Root struct {
	A    *int
	B    *float64
	T    *time.Time
	Z    *Zonk
	NaNa *int
	L    []string
}

func Test_lab(t *testing.T) {
	pp := &Root{}
	require.NoError(t, Fill(pp))

}
