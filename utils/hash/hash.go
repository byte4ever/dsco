package hash

import (
	"crypto"
	"errors"
	"fmt"

	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/md4"
	_ "golang.org/x/crypto/ripemd160"
	_ "golang.org/x/crypto/sha3"
)

// ErrInvalidHash when Hash cannot be cast properly.
var ErrInvalidHash = errors.New("invalid hash")

// Hash is a type to ensure crypto.Hash() can  support various marshaling techniques.
type Hash crypto.Hash

func (d *Hash) UnmarshalText(text []byte) error {
	s := string(text)
	h, found := tab[s]

	if !found {
		return fmt.Errorf("unknown hash %s: %w", s, ErrInvalidHash)
	}

	*d = Hash(h)

	return nil
}

func (d *Hash) MarshalText() (text []byte, err error) {
	return []byte(crypto.Hash(*d).String()), nil
}

func (d *Hash) ToCrypto() crypto.Hash {
	return crypto.Hash(*d)
}

var tab = map[string]crypto.Hash{
	"MD4":         crypto.MD4,
	"MD5":         crypto.MD5,
	"SHA-1":       crypto.SHA1,
	"SHA-224":     crypto.SHA224,
	"SHA-256":     crypto.SHA256,
	"SHA-384":     crypto.SHA384,
	"SHA-512":     crypto.SHA512,
	"RIPEMD-160":  crypto.RIPEMD160,
	"SHA3-224":    crypto.SHA3_224,
	"SHA3-256":    crypto.SHA3_256,
	"SHA3-384":    crypto.SHA3_384,
	"SHA3-512":    crypto.SHA3_512,
	"SHA-512/224": crypto.SHA512_224,
	"SHA-512/256": crypto.SHA512_256,
	"BLAKE2s-256": crypto.BLAKE2s_256,
	"BLAKE2b-256": crypto.BLAKE2b_256,
	"BLAKE2b-384": crypto.BLAKE2b_384,
	"BLAKE2b-512": crypto.BLAKE2b_512,
}
