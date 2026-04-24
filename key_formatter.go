package dsco

import "strings"

type (
	// KeyFormatter renders a layer-internal alias path into the canonical
	// user-facing key form for a specific layer kind (env, cmdline, file).
	//
	// Pattern: Strategy — each layer kind formats keys differently; the
	// formatter is injected into StringBasedBuilder at construction time.
	KeyFormatter interface { //nolint:iface // strategy interface; consumed by StringBasedBuilder in task 5
		// LayerKind returns the layer category (e.g. "env", "cmdline", "file").
		// Empty for layers that cannot enumerate keys.
		LayerKind() string

		// LayerName returns the layer instance identifier
		// (e.g. "env:MYAPP", "cmdline", "file:config.yaml").
		LayerName() string

		// FormatKey converts an internal alias path (dash-separated, lowercase)
		// to the canonical user-facing key for this layer kind. Returns the
		// empty string when the layer cannot enumerate keys.
		FormatKey(aliasPath string) string
	}

	// envKeyFormatter formats keys for environment-variable layers:
	// PREFIX-UPPER-CASE-DASHED.
	envKeyFormatter struct { //nolint:unused // consumed by StringBasedBuilder in task 5
		prefix string
	}

	// cmdlineKeyFormatter formats keys for command-line layers: --name=.
	cmdlineKeyFormatter struct{} //nolint:unused // consumed by StringBasedBuilder in task 5

	// fileKeyFormatter formats keys for file layers: dot-separated YAML path.
	fileKeyFormatter struct { //nolint:unused // consumed by StringBasedBuilder in task 5
		id string
	}

	// nilKeyFormatter is a no-op formatter for layers (custom string
	// providers) that cannot enumerate keys statically. LayerKind is empty so
	// reduce-pass logic skips them when picking a canonical key.
	nilKeyFormatter struct { //nolint:unused // consumed by StringBasedBuilder in task 5
		name string
	}
)

func newEnvKeyFormatter(prefix string) *envKeyFormatter { //nolint:unused // consumed by StringBasedBuilder in task 5
	return &envKeyFormatter{prefix: prefix}
}

func newCmdlineKeyFormatter() *cmdlineKeyFormatter { //nolint:unused // consumed by StringBasedBuilder in task 5
	return &cmdlineKeyFormatter{}
}

func newFileKeyFormatter(id string) *fileKeyFormatter { //nolint:unused // consumed by StringBasedBuilder in task 5
	return &fileKeyFormatter{id: id}
}

func newNilKeyFormatter(name string) *nilKeyFormatter { //nolint:unused // consumed by StringBasedBuilder in task 5
	return &nilKeyFormatter{name: name}
}

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*envKeyFormatter) LayerKind() string { return "env" }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (f *envKeyFormatter) LayerName() string { return "env:" + f.prefix }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (f *envKeyFormatter) FormatKey(aliasPath string) string {
	return f.prefix + "-" + strings.ToUpper(aliasPath)
}

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*cmdlineKeyFormatter) LayerKind() string { return "cmdline" }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*cmdlineKeyFormatter) LayerName() string { return "cmdline" }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*cmdlineKeyFormatter) FormatKey(aliasPath string) string {
	return "--" + aliasPath + "="
}

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*fileKeyFormatter) LayerKind() string { return "file" }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (f *fileKeyFormatter) LayerName() string { return "file:" + f.id }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*fileKeyFormatter) FormatKey(aliasPath string) string {
	return strings.ReplaceAll(aliasPath, "-", ".")
}

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*nilKeyFormatter) LayerKind() string { return "" }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (f *nilKeyFormatter) LayerName() string { return f.name }

//nolint:unused // consumed by StringBasedBuilder in task 5
func (*nilKeyFormatter) FormatKey(_ string) string { return "" }
