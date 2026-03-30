package hprovider

import (
	"hash"
)

// Provider manages a pool of hash instances for efficient reuse.
// It uses a channel-based pool to provide and recycle hash instances,
// reducing allocation overhead in hash-intensive operations.
type Provider[T hash.Hash] struct {
	provideFunc func() T // Factory function for creating new hash instances
	store       chan T   // Channel-based pool for storing reusable instances
}

// New creates a new hash provider with the specified factory function and pool
// size.
//
// Parameters:
// - provideFunc: Factory function that creates new hash instances when the pool
// is empty
//   - storeSize: Maximum number of hash instances to keep in the pool for reuse
//
// Returns a new Provider that manages hash instance pooling.
func New[T hash.Hash](
	provideFunc func() T,
	storeSize int,
) *Provider[T] {
	return &Provider[T]{
		provideFunc: provideFunc,
		store:       make(chan T, storeSize),
	}
}

// Get retrieves a hash instance from the pool or creates a new one.
// If the pool has available instances, one is retrieved, reset, and returned.
// If the pool is empty, a new instance is created using the factory function.
// This method is non-blocking and thread-safe.
func (h *Provider[T]) Get() T {
	select {
	case v := <-h.store:
		v.Reset()
		return v
	default:
		return h.provideFunc()
	}
}

// PutBack returns a hash instance to the pool for future reuse.
// If the pool has space, the instance is stored for later use.
// If the pool is full, the instance is discarded and will be garbage collected.
// This method is non-blocking and thread-safe.
func (h *Provider[T]) PutBack(v T) {
	select {
	case h.store <- v:
		return
	default:
		return
	}
}
