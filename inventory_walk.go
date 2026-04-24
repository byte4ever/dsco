package dsco

import (
	"errors"
	"fmt"
)

// InventoryWalk holds the prepared model and per-layer reporters used by
// the inventory sub-package to compute a Report. It contains no live
// configuration values — only structural metadata.
type InventoryWalk struct {
	Model     ModelInterface
	Reporters []InventoryReporter
}

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
