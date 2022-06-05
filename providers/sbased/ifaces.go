package sbased

import (
	"github.com/byte4ever/dsco"
)

// EntriesProvider defines string based entries provider behaviour.
type EntriesProvider interface {

	// GetOrigin returns the origin id of the provider.
	GetOrigin() dsco.Origin

	// GetEntries returns the string based entries.
	GetEntries() Entries
}

type entry struct {
	Entry
	bounded bool
	used    bool
}

type entries map[string]*entry

// Entry is single entry for string based binder.
type Entry struct {
	ExternalKey string // is the external key.
	Value       string // is the value.
}

// Entries is a Key/Entry set.
type Entries map[string]*Entry
