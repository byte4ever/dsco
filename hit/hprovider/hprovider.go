package hprovider

import (
	"hash"
)

type HashProviderImpl[T hash.Hash] struct {
	provideFunc func() T
	store       chan T
}

func NewHashProviderImpl[T hash.Hash](
	provideFunc func() T,
	storeSize int,
) *HashProviderImpl[T] {
	return &HashProviderImpl[T]{
		provideFunc: provideFunc,
		store:       make(chan T, storeSize),
	}
}
