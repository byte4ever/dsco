package dsco

type (
	// InventoryReporter is implemented by layer builders that can describe
	// what they would contribute to an inventory without performing I/O.
	//
	// Pattern: Strategy — each layer reports its own contribution shape (key
	// form for string-based layers, baked-in values for struct layers).
	InventoryReporter interface { //nolint:iface // public API; used by external implementors
		// ReportInventory returns the layer's contribution to an inventory
		// walk for the given model. Implementations must not perform any I/O.
		ReportInventory(model ModelInterface) (LayerInventory, error)
	}

	// LayerInventory is one layer's contribution to a Report.
	LayerInventory struct {
		// Name uniquely identifies the layer instance, e.g. "env:MYAPP",
		// "cmdline", "file:<id>", "struct:<id>", or a custom provider name.
		Name string

		// Optional information for callers that cannot enumerate keys
		// (typically custom string providers).
		Note string

		// Provides lists every (field, key|value) pair this layer can
		// supply to the model.
		Provides []FieldProvision
	}

	// FieldProvision is one (field, layer) pair from a layer's perspective.
	FieldProvision struct {
		// Value is the baked-in value for struct layers; nil for
		// string-based layers.
		Value any

		// FieldUID matches the model's field UID.
		FieldUID string

		// Key is the canonical key for string-based layers; empty for
		// struct layers.
		Key string
	}
)
