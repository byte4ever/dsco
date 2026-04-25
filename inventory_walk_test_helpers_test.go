package dsco

// This file contains test-only layer helpers for exercising
// PrepareInventoryWalk and Compute error paths. These types are compiled only
// during `go test` and are invisible to package consumers.

import "github.com/byte4ever/dsco/internal/fvalue"

type (
	// testBareFieldValuesGetter is a FieldValuesGetter that does NOT
	// implement InventoryReporter. Used to exercise the
	// ErrLayerNotInventoryReporter branch in PrepareInventoryWalk.
	testBareFieldValuesGetter struct{}

	// testNonReporterLayer is a Layer whose FieldValuesGetter does not
	// implement InventoryReporter.
	testNonReporterLayer struct{}

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
)

func (testBareFieldValuesGetter) GetFieldValuesFrom(
	_ ModelInterface,
) (fvalue.Values, error) {
	return nil, nil //nolint:nilnil // test stub
}

func (testNonReporterLayer) register(to *layerBuilder) error {
	to.addBuilder(newNormalLayer(testBareFieldValuesGetter{}))
	return nil
}

// WithNonReporterLayerForTest returns a Layer whose FieldValuesGetter does
// not implement InventoryReporter. Used solely by tests to exercise the
// ErrLayerNotInventoryReporter branch in PrepareInventoryWalk.
//
//nolint:ireturn,iface // test seam; wraps concrete type as Layer interface
func WithNonReporterLayerForTest() Layer {
	return testNonReporterLayer{}
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
