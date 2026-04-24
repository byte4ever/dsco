package dsco

import (
	"errors"
	"fmt"

	"github.com/byte4ever/dsco/internal/fvalue"
)

type (
	// minimalFieldValuesGetter is a bare-bones FieldValuesGetter that does
	// NOT implement InventoryReporter. Used only by
	// WithNonReporterLayerForTest.
	minimalFieldValuesGetter struct{}

	// nonReporterLayer is a test-only Layer whose underlying
	// FieldValuesGetter does not implement InventoryReporter. Used to
	// exercise the ErrLayerNotInventoryReporter branch in
	// PrepareInventoryWalk.
	nonReporterLayer struct{}

	// inventoryReporterFVG adapts an InventoryReporter to the
	// FieldValuesGetter interface, allowing it to act as a layer inside
	// the policy pipeline.
	inventoryReporterFVG struct {
		InventoryReporter
	}

	// inventoryReporterLayer wraps an InventoryReporter as a full Layer so
	// tests can inject synthetic reporters into
	// PrepareInventoryWalk/Compute.
	inventoryReporterLayer struct {
		reporter InventoryReporter
	}

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

func (minimalFieldValuesGetter) GetFieldValuesFrom(
	_ ModelInterface,
) (fvalue.Values, error) {
	return nil, nil //nolint:nilnil // test stub — never called in normal flow
}

func (nonReporterLayer) register(to *layerBuilder) error {
	to.addBuilder(newNormalLayer(minimalFieldValuesGetter{}))
	return nil
}

// WithNonReporterLayerForTest returns a Layer whose FieldValuesGetter does
// not implement InventoryReporter. Intended solely for tests that need to
// trigger the ErrLayerNotInventoryReporter branch in PrepareInventoryWalk.
//
//nolint:ireturn,iface // test seam; wraps concrete type as Layer interface
func WithNonReporterLayerForTest() Layer {
	return nonReporterLayer{}
}

func (*inventoryReporterFVG) GetFieldValuesFrom(
	_ ModelInterface,
) (fvalue.Values, error) {
	return nil, nil //nolint:nilnil // synthetic layer — values not needed
}

func (l inventoryReporterLayer) register(to *layerBuilder) error {
	to.addBuilder(newNormalLayer(&inventoryReporterFVG{l.reporter}))
	return nil
}

// WithInventoryReporterLayerForTest wraps reporter as a normal Layer.
// Intended solely for tests that need to inject a synthetic
// InventoryReporter into Compute to exercise the ReportInventory error
// path.
//
//nolint:ireturn,iface // test seam; wraps concrete type as Layer interface
func WithInventoryReporterLayerForTest(reporter InventoryReporter) Layer {
	return inventoryReporterLayer{reporter: reporter}
}

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
