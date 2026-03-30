package hit

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/byte4ever/dsco/hit/hprovider"
)

func TestNewIntNode(t *testing.T) {
	var nodes []MerkelNode

	hp := hprovider.New(md5.New, 100)

	for i := 0; i < 10; i++ {
		k1 := newkeyedNode(
			hp,
			nil,
			"k1",
			NewIntNode(
				hp,
				nil,
				"polo",
				i,
			),
			NewStringNode(
				hp,
				nil,
				"s1",
				fmt.Sprintf("lola %d", i),
			),
		)
		nodes = append(nodes, k1)
	}

	p := newParentNode(
		hp,
		nil,
		"parent",
		nodes...,
	)

	fmt.Println(hex.EncodeToString(p.GetHash()))
}

func TestNewStringNode(t *testing.T) {
	var nodes []MerkelNode

	hp := hprovider.New(sha1.New, 100)

	m := map[string]string{
		"rose":    "martin-bey",
		"lola":    "martin-bey",
		"celine":  "bey",
		"laurent": "martin",
	}

	for key, val := range m {
		k1 := newkeyedNode(
			hp,
			nil,
			"k1",
			NewStringNode(
				hp,
				nil,
				"",
				key,
			),
			NewStringNode(
				hp,
				nil,
				"",
				val,
			),
		)
		nodes = append(nodes, k1)
	}

	p := newParentNode(
		hp,
		nil,
		"parent",
		nodes...,
	)

	fmt.Println(hex.EncodeToString(p.GetHash()))
}

func TestGetID(t *testing.T) {
	hp := hprovider.New(md5.New, 10)

	// Test nodeImpl GetID method.
	n := &nodeImpl{
		id:   "test-node-id",
		hash: []byte{1, 2, 3, 4},
	}

	if n.GetID() != "test-node-id" {
		t.Errorf("Expected ID 'test-node-id', got '%s'", n.GetID())
	}

	// Test IntNode GetID method.
	intNode := NewIntNode(
		hp,
		nil,
		"int-node-123",
		42,
	)
	if intNode.GetID() != "int-node-123" {
		t.Errorf("Expected ID 'int-node-123', got '%s'", intNode.GetID())
	}

	// Test StringNode GetID method.
	stringNode := NewStringNode(
		hp,
		nil,
		"string-node-456",
		"test-value",
	)
	if stringNode.GetID() != "string-node-456" {
		t.Errorf("Expected ID 'string-node-456', got '%s'", stringNode.GetID())
	}

	// Test keyedNode GetID method.
	keyNode := newkeyedNode(
		hp,
		nil,
		"keyed-node-789",
		NewStringNode(
			hp,
			nil,
			"key",
			"test-key",
		),
		NewStringNode(
			hp,
			nil,
			"value",
			"test-value",
		),
	)
	if keyNode.GetID() != "keyed-node-789" {
		t.Errorf("Expected ID 'keyed-node-789', got '%s'", keyNode.GetID())
	}

	// Test parentNode GetID method.
	parentNode := newParentNode(
		hp,
		nil,
		"parent-node-000",
		intNode,
		stringNode,
	)
	if parentNode.GetID() != "parent-node-000" {
		t.Errorf("Expected ID 'parent-node-000', got '%s'", parentNode.GetID())
	}
}
