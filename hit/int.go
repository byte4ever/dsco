package hit

import (
	"crypto/subtle"
	"encoding/binary"
	"hash"
)

var (
	saltInt      = []byte("int-d4a64668-b28d-4ed2-9760-3f498e663c5e")
	saltString   = []byte("string-3981a84a-fed5-44e8-8093-b128c84a377b")
	saltInteKeyd = []byte("intKeyed-5b71742e-7a62-4170-adde-cf8140babcc1")
)

// NodeID represents a unique identifier for hash nodes in the Merkle-like tree
// structure. It provides a string-based identifier that can be used to
// reference and locate
// specific nodes within the hash tree.
type NodeID string

// HashProvider defines an interface for providing hash instances with pooling
// support. This generic interface allows efficient reuse of hash instances,
// reducing allocation
// overhead during repeated hash computations.
type HashProvider[T hash.Hash] interface {
	// Get retrieves a hash instance from the pool or creates a new one.
	// The returned instance is ready for use and should be returned via PutBack
	// when done.
	Get() T

	// PutBack returns a hash instance to the pool for reuse.
	// This enables efficient resource management and reduces garbage collection
	// pressure.
	PutBack(h T)
}

// MerkelNode defines the interface for nodes in a Merkle-like hash tree
// structure.
// Each node can provide its unique identifier and computed hash value, enabling
// efficient tree traversal and verification.
type MerkelNode interface {
	// GetID returns the unique identifier for this node.
	// The ID can be used to reference and locate the node within the tree
	// structure.
	GetID() string

	// GetHash returns the computed hash value for this node.
	// The hash represents the content and structure of the node and its
	// children.
	GetHash() []byte
}

type nodeImpl struct {
	id   string
	hash []byte
}

func (n *nodeImpl) GetID() string {
	return n.id
}

func (n *nodeImpl) GetHash() []byte {
	return n.hash
}

// IntNode represents a hash tree node that contains an integer value.
// It implements the MerkelNode interface and provides content-addressable
// hashing for integer data.
type IntNode struct {
	nodeImpl     // Embedded base implementation
	value    int // The integer value stored in this node
}

// NewIntNode creates a new IntNode with the specified value and computes its
// hash. The hash is computed using the provided salt for security and collision
// resistance.
//
// Parameters:
//   - hashProvider: Provider for hash instances to avoid allocation overhead
//   - salt: Random bytes to prevent hash collision attacks
//   - id: Unique identifier for this node
//   - value: The integer value to store and hash
//
// Returns a new IntNode with computed hash value.
func NewIntNode(
	hashProvider HashProvider[hash.Hash],
	salt []byte,
	id string,
	value int,
) *IntNode {
	buf := make([]byte, binary.MaxVarintLen64)

	h := hashProvider.Get()
	defer hashProvider.PutBack(h)

	sig := make(
		[]byte,
		0,
		h.Size(),
	)

	h.Write(salt)
	h.Write(saltInt)
	h.Write(salt)
	binary.PutVarint(buf, int64(value))
	h.Write(salt)
	h.Write(buf)
	h.Write(salt)
	h.Write(saltInt)
	h.Write(salt)

	sig = h.Sum(sig)

	return &IntNode{
		value: value,
		nodeImpl: nodeImpl{
			id:   id,
			hash: sig,
		},
	}
}

type parentNode struct {
	nodeImpl
	children []MerkelNode
}

func newParentNode(
	hashProvider HashProvider[hash.Hash],
	salt []byte,
	id string,
	node ...MerkelNode,
) *parentNode {
	h := hashProvider.Get()
	defer hashProvider.PutBack(h)

	accumSig := make(
		[]byte,
		0,
		h.Size(),
	)
	accum := make([]byte, h.Size())

	saltCount := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(saltCount, uint64(len(node)))

	for _, merkelNode := range node {
		subtle.XORBytes(
			accum,
			accum,
			merkelNode.GetHash(),
		)
	}

	h.Write(salt)
	h.Write(saltCount)
	h.Write(accum)
	h.Write(saltCount)
	h.Write(salt)

	accumSig = h.Sum(accumSig)

	return &parentNode{
		nodeImpl: nodeImpl{
			id:   id,
			hash: accumSig,
		},
		children: node,
	}
}

// StringNode represents a hash tree node that contains a string value.
// It implements the MerkelNode interface and provides content-addressable
// hashing for string data.
type StringNode struct {
	value    string // The string value stored in this node
	nodeImpl        // Embedded base implementation
}

// NewStringNode creates a new StringNode with the specified value and computes
// its hash. The hash is computed using the provided salt for security and
// collision resistance.
//
// Parameters:
//   - hashProvider: Provider for hash instances to avoid allocation overhead
//   - salt: Random bytes to prevent hash collision attacks
//   - id: Unique identifier for this node
//   - value: The string value to store and hash
//
// Returns a new StringNode with computed hash value.
func NewStringNode(
	hashProvider HashProvider[hash.Hash],
	salt []byte,
	id string,
	value string,
) *StringNode {
	buf := []byte(value)

	h := hashProvider.Get()
	defer hashProvider.PutBack(h)

	sig := make(
		[]byte,
		0,
		h.Size(),
	)

	h.Write(salt)
	h.Write(saltString)
	h.Write(salt)
	h.Write(buf)
	h.Write(salt)
	h.Write(saltString)
	h.Write(salt)

	sig = h.Sum(sig)

	return &StringNode{
		value: value,
		nodeImpl: nodeImpl{
			id:   id,
			hash: sig,
		},
	}
}

type keyedNode struct {
	key   MerkelNode
	value MerkelNode
	nodeImpl
}

func newkeyedNode(
	hashProvider HashProvider[hash.Hash],
	salt []byte,
	id string,
	key MerkelNode,
	value MerkelNode,
) *keyedNode {
	h := hashProvider.Get()
	defer hashProvider.PutBack(h)

	hashSig := make(
		[]byte,
		0,
		h.Size(),
	)

	h.Write(salt)
	h.Write(saltInteKeyd)
	h.Write(salt)
	h.Write(key.GetHash())
	h.Write(salt)
	h.Write(saltInteKeyd)
	h.Write(salt)
	h.Write(value.GetHash())
	h.Write(salt)
	h.Write(saltInteKeyd)
	h.Write(salt)

	hashSig = h.Sum(hashSig)

	return &keyedNode{
		key:   key,
		value: value,
		nodeImpl: nodeImpl{
			id:   id,
			hash: hashSig,
		},
	}
}
