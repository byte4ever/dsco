package kfile

import (
	"errors"
	"testing"
)

var errMocked1 = errors.New("err1")
var errMocked2 = errors.New("err2")

func TestPathErrors_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
		p    PathErrors
	}{
		{
			name: "t1",
			p: PathErrors{
				&pathError{
					err:  errMocked1,
					path: "p1",
				},
				&pathError{
					err:  errMocked2,
					path: "p2",
				},
			},
			want: "p1: err1\np2: err2\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				if got := tt.p.Error(); got != tt.want {
					t.Errorf("Error() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
