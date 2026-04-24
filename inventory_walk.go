package dsco

import (
	"errors"
	"fmt"
)

type (
	// InventoryWalk holds the prepared model and per-layer reporters used
	// by the inventory sub-package to compute a Report. It contains no
	// live configuration values — only structural metadata.
	InventoryWalk struct {
		Model     ModelInterface
		Reporters []InventoryReporter
	}
)

// ErrLayerNotInventoryReporter indicates that a layer's underlying
// FieldValuesGetter does not implement InventoryReporter, which is
// required for PrepareInventoryWalk.
var ErrLayerNotInventoryReporter = errors.New(
	"layer does not implement InventoryReporter",
)

// BuildModel constructs the configuration model from a pointer-to-struct
// value, mirroring the model-build phase of Fill. Exposed for the
// inventory sub-package; callers should generally use Fill or
// PrepareInventoryWalk instead.
//
//nolint:iface,ireturn,revive // returns shared ModelInterface; name intentionally mirrors buildModel
func BuildModel(cfg any) (ModelInterface, error) {
	return buildModel(cfg) //nolint:wrapcheck // thin public wrapper; caller sees buildModel errors directly
}

// CollectAliasesForTest exposes collectAliases for testing only. It is
// intentionally not part of the public API.
func CollectAliasesForTest(mdl ModelInterface) (map[string]string, error) {
	return collectAliases(mdl) //nolint:wrapcheck // thin test seam; caller sees collectAliases errors directly
}

// prepareInventoryWalkFromPolicies is the internal implementation shared by
// PrepareInventoryWalk and in-package tests. It converts a
// constraintLayerPolicies slice into an InventoryWalk, verifying that every
// policy's FieldValuesGetter implements InventoryReporter.
func prepareInventoryWalkFromPolicies(
	policies constraintLayerPolicies,
	mdl ModelInterface,
) (*InventoryWalk, error) {
	const errCtx = "preparing inventory walk"

	reporters := make([]InventoryReporter, 0, len(policies))

	for i, p := range policies {
		fvg := p.getFieldValuesGetter()

		reporter, ok := fvg.(InventoryReporter)
		if !ok {
			return nil, fmt.Errorf(
				"%s: layer #%d: %w", errCtx, i, ErrLayerNotInventoryReporter,
			)
		}

		reporters = append(reporters, reporter)
	}

	return &InventoryWalk{
		Model:     mdl,
		Reporters: reporters,
	}, nil
}

// PrepareInventoryWalk constructs the model and the per-layer
// InventoryReporter list without performing any I/O. Used by the
// inventory sub-package to compute a Report.
//
// Pattern: Factory — assembles the structural inputs needed for an
// inventory walk.
func PrepareInventoryWalk(
	cfg any,
	layers ...Layer,
) (*InventoryWalk, error) {
	const errCtx = "preparing inventory walk"

	mdl, err := buildModel(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	policies, err := Layers(layers).GetPolicies()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	walk, err := prepareInventoryWalkFromPolicies(policies, mdl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errCtx, err)
	}

	return walk, nil
}
