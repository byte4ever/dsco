package sbased

import (
	"github.com/byte4ever/dsco"
)

type StrEntriesProvider interface {
	GetOrigin() dsco.Origin
	GetEntries() StrEntries
}

type entry struct {
	StrEntry
	bounded bool
	used    bool
}

type entries map[string]*entry

type StrEntry struct {
	ExternalKey string
	Value       string
}

type StrEntries map[string]*StrEntry
