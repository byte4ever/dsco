package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type FakeStruct struct {
	H *Hash `yaml:"h2"`
}

func hp(h Hash) *Hash {
	v := h
	return &v
}

func TestHash_MarshalYAML(t *testing.T) {
	t.Run(
		"all succeed", func(t *testing.T) {
			for _, hash := range tab {
				fori := &FakeStruct{
					H: hp(Hash(hash)),
				}
				bb, err := yaml.Marshal(fori)
				require.Nil(t, err)
				require.NotNil(t, bb)

				var fs FakeStruct
				err = yaml.Unmarshal(bb, &fs)
				require.Nil(t, err)
				require.Equal(t, fori.H.ToCrypto(), fs.H.ToCrypto())
				require.True(t, fs.H.ToCrypto().Available())
			}
		},
	)

	t.Run(
		"unmarshal test failure", func(t *testing.T) {
			var h Hash
			err := h.UnmarshalText([]byte("INVALID_HASH"))
			require.ErrorIs(t, err, ErrInvalidHash)
			require.ErrorContains(t, err, "INVALID_HASH")
		},
	)
}
