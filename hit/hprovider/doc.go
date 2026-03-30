/*
Package hprovider provides hash instance pooling for efficient hash computation
in dsco's hit package.

# Overview

The hprovider package implements a generic hash provider that manages a pool
of hash instances for reuse. This reduces allocation overhead when performing
many hash computations, which is common in dsco's configuration processing
where hash values are computed for deduplication, caching, and content
addressing.

# Core Type

## Provider

Provider is a generic type that manages a pool of hash instances:

	type Provider[T hash.Hash] struct {
		provideFunc func() T    // Factory function for new hash instances
		store       chan T      // Channel-based pool for reusing instances
	}

The Provider type uses Go generics to work with any hash.Hash implementation,
providing type safety while maintaining efficiency.

# Constructor

## New Function

New creates a new hash provider with specified factory function and pool size:

	func New[T hash.Hash](provideFunc func() T, storeSize int) *Provider[T]

### Parameters

  - **provideFunc**: A factory function that creates new hash instances. This
    is called when the pool is empty and a new hash instance is needed.

  - **storeSize**: The maximum number of hash instances to keep in the pool.
    A larger pool reduces allocations but uses more memory.

### Examples

	// MD5 hash provider
	md5Provider := hprovider.New(md5.New, 10)

	// SHA256 hash provider
	sha256Provider := hprovider.New(sha256.New, 5)

	// Custom hash provider
	customProvider := hprovider.New(func() *CustomHash {
		return &CustomHash{} // Custom initialization
	}, 20)

# Pool Management

## Get Method

Get retrieves a hash instance from the pool or creates a new one:

	func (h *Provider[T]) Get() T

### Behavior

1. **Pool hit**: If an instance is available in the pool, it's retrieved and reset
2. **Pool miss**: If the pool is empty, a new instance is created using the factory function
3. **Non-blocking**: Uses a select statement to avoid blocking on empty pool

### Usage

	hashInstance := provider.Get()
	defer provider.PutBack(hashInstance)

	// Use the hash instance
	hashInstance.Write(data)
	result := hashInstance.Sum(nil)

## PutBack Method

PutBack returns a hash instance to the pool for reuse:

	func (h *Provider[T]) PutBack(v T)

### Behavior

1. **Pool space available**: Instance is stored in the pool for reuse
2. **Pool full**: Instance is discarded (garbage collected)
3. **Non-blocking**: Uses a select statement to avoid blocking on full pool

### Usage

	hashInstance := provider.Get()

	// Use the hash instance
	hashInstance.Write(data)
	result := hashInstance.Sum(nil)

	// Return to pool when done
	provider.PutBack(hashInstance)

# Usage Patterns

## Basic Usage Pattern

The typical usage pattern with defer for automatic cleanup:

	func computeHash(provider *hprovider.Provider[hash.Hash], data []byte) []byte {
		h := provider.Get()
		defer provider.PutBack(h)

		h.Write(data)
		return h.Sum(nil)
	}

## Batch Processing Pattern

For processing multiple items efficiently:

	func computeHashes(provider *hprovider.Provider[hash.Hash], items [][]byte) [][]byte {
		results := make([][]byte, len(items))

		h := provider.Get()
		defer provider.PutBack(h)

		for i, item := range items {
			h.Reset()  // Reset for each item
			h.Write(item)
			results[i] = h.Sum(nil)
		}

		return results
	}

## Concurrent Usage Pattern

Provider instances are safe for concurrent use:

	func processItemsConcurrently(provider *hprovider.Provider[hash.Hash], items [][]byte) [][]byte {
		results := make([][]byte, len(items))
		var wg sync.WaitGroup

		for i, item := range items {
			wg.Add(1)
			go func(index int, data []byte) {
				defer wg.Done()

				h := provider.Get()
				defer provider.PutBack(h)

				h.Write(data)
				results[index] = h.Sum(nil)
			}(i, item)
		}

		wg.Wait()
		return results
	}

# Integration with dsco Components

## Node Creation

Used in hit package for efficient hash computation during node creation:

	type IntNode struct {
		nodeImpl
		value int
	}

	func NewIntNode(
		hashProvider hprovider.Provider[hash.Hash],
		salt []byte,
		id string,
		value int,
	) *IntNode {
		h := hashProvider.Get()
		defer hashProvider.PutBack(h)

		// Compute hash with salt and value
		h.Write(salt)
		binary.PutVarint(buf, int64(value))
		h.Write(buf)

		return &IntNode{
			value: value,
			nodeImpl: nodeImpl{
				id:   id,
				hash: h.Sum(nil),
			},
		}
	}

## Configuration Hashing

Used for computing configuration hashes in dsco:

	func hashConfiguration(provider *hprovider.Provider[hash.Hash], config interface{}) uint64 {
		h := provider.Get()
		defer provider.PutBack(h)

		// Serialize and hash configuration
		data, _ := json.Marshal(config)
		h.Write(data)
		hashBytes := h.Sum(nil)

		// Convert to uint64 for use as map key
		return binary.BigEndian.Uint64(hashBytes[:8])
	}

# Performance Characteristics

## Pool Efficiency

The channel-based pool provides excellent performance:

- **O(1) Get operation**: Channel receive is constant time
- **O(1) PutBack operation**: Channel send is constant time
- **Non-blocking operations**: No goroutine blocking on pool operations
- **Automatic overflow handling**: Full pool discards instances gracefully

## Memory Management

Pool size affects memory vs allocation trade-offs:

- **Small pools (1-5)**: Lower memory usage, more allocations
- **Medium pools (10-50)**: Balanced memory and allocation efficiency
- **Large pools (100+)**: Higher memory usage, minimal allocations

## Allocation Reduction

Effective pooling can reduce allocations significantly:

	// Without pooling: Creates new hash for each operation
	for i := 0; i < 1000; i++ {
		h := sha256.New()  // 1000 allocations
		h.Write(data[i])
		results[i] = h.Sum(nil)
	}

	// With pooling: Reuses hash instances
	provider := hprovider.New(sha256.New, 10)
	for i := 0; i < 1000; i++ {
		h := provider.Get()  // ~10 allocations total
		h.Write(data[i])
		results[i] = h.Sum(nil)
		provider.PutBack(h)
	}

# Thread Safety

The Provider type is designed for concurrent use:

- **Get() is thread-safe**: Multiple goroutines can call Get() concurrently
- **PutBack() is thread-safe**: Multiple goroutines can call PutBack() concurrently
- **Channel synchronization**: Uses Go channels for thread-safe pool operations
- **No external locking required**: All synchronization is internal

# Best Practices

## Pool Sizing

Choose pool size based on concurrency and memory constraints:

	// High concurrency, memory available
	provider := hprovider.New(sha256.New, 50)

	// Low concurrency, memory constrained
	provider := hprovider.New(sha256.New, 5)

	// Single-threaded usage
	provider := hprovider.New(sha256.New, 1)

## Resource Management

Always use defer for proper resource management:

	// Good: Guaranteed cleanup
	func useHash(provider *hprovider.Provider[hash.Hash]) []byte {
		h := provider.Get()
		defer provider.PutBack(h)  // Always returned to pool

		// ... use hash ...
		return result
	}

	// Risky: Manual cleanup might be missed
	func useHash(provider *hprovider.Provider[hash.Hash]) []byte {
		h := provider.Get()

		// ... use hash ...

		provider.PutBack(h)  // Might be skipped on early return
		return result
	}

## Hash Instance Reuse

Reset hash instances when reusing in loops:

	h := provider.Get()
	defer provider.PutBack(h)

	for _, item := range items {
		h.Reset()  // Clear previous state
		h.Write(item)
		results = append(results, h.Sum(nil))
	}

# Error Handling

Provider operations are designed to be robust:

- **Get() never fails**: Always returns a valid hash instance
- **PutBack() never fails**: Discards instance if pool is full
- **No error returns**: Simplified API with predictable behavior

# Memory Considerations

The Provider manages memory efficiently:

- **Bounded memory usage**: Pool size limits maximum memory usage
- **Garbage collection friendly**: Unused instances are collected normally
- **No memory leaks**: Full pools automatically discard excess instances

# Integration Examples

## Complete Hash Processing Pipeline

	type HashingService struct {
		md5Provider    *hprovider.Provider[hash.Hash]
		sha256Provider *hprovider.Provider[hash.Hash]
	}

	func NewHashingService() *HashingService {
		return &HashingService{
			md5Provider:    hprovider.New(md5.New, 10),
			sha256Provider: hprovider.New(sha256.New, 10),
		}
	}

	func (s *HashingService) ComputeHashes(data []byte) (md5Hash, sha256Hash []byte) {
		// Compute MD5
		md5h := s.md5Provider.Get()
		defer s.md5Provider.PutBack(md5h)
		md5h.Write(data)
		md5Hash = md5h.Sum(nil)

		// Compute SHA256
		sha256h := s.sha256Provider.Get()
		defer s.sha256Provider.PutBack(sha256h)
		sha256h.Write(data)
		sha256Hash = sha256h.Sum(nil)

		return md5Hash, sha256Hash
	}

The hprovider package provides essential hash instance pooling that enables
efficient hash computation throughout dsco's configuration processing,
reducing allocation overhead and improving overall performance.
*/
package hprovider
