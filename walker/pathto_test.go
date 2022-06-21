package walker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pathTo(t *testing.T) {
	t.Parallel()

	type args struct {
		currentPath string
		fieldName   string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no current path",
			args: args{
				currentPath: "",
				fieldName:   "field",
			},
			want: "field",
		},
		{
			name: "with current path",
			args: args{
				currentPath: "curPath",
				fieldName:   "field",
			},
			want: "curPath.field",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(
			tt.name,
			func(t *testing.T) {
				t.Parallel()

				require.Equal(
					t,
					tt.want,
					pathTo(tt.args.currentPath, tt.args.fieldName),
				)
			},
		)
	}
}
