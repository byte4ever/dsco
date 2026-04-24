package dsco

import "strings"

type (
	// KeyFormatter renders a layer-internal alias path into the canonical
	// user-facing key form for a specific layer kind (env, cmdline, file).
	//
	// Pattern: Strategy — each layer kind formats keys differently; the
	// formatter is injected into StringBasedBuilder at construction time.
	KeyFormatter interface {
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
	envKeyFormatter struct {
		prefix string
	}

	// cmdlineKeyFormatter formats keys for command-line layers: --name=.
	cmdlineKeyFormatter struct{}

	// nilKeyFormatter is a no-op formatter for layers (custom string
	// providers) that cannot enumerate keys statically. LayerKind is empty so
	// reduce-pass logic skips them when picking a canonical key.
	nilKeyFormatter struct {
		name string
	}
)

func newEnvKeyFormatter(prefix string) *envKeyFormatter {
	return &envKeyFormatter{prefix: prefix}
}

func newCmdlineKeyFormatter() *cmdlineKeyFormatter {
	return &cmdlineKeyFormatter{}
}

func newNilKeyFormatter(name string) *nilKeyFormatter {
	return &nilKeyFormatter{name: name}
}

func (*envKeyFormatter) LayerKind() string { return "env" }

func (f *envKeyFormatter) LayerName() string { return "env:" + f.prefix }

func (f *envKeyFormatter) FormatKey(aliasPath string) string {
	return f.prefix + "-" + strings.ToUpper(aliasPath)
}

func (*cmdlineKeyFormatter) LayerKind() string { return "cmdline" }

func (*cmdlineKeyFormatter) LayerName() string { return "cmdline" }

func (*cmdlineKeyFormatter) FormatKey(aliasPath string) string {
	return "--" + aliasPath + "="
}

func (*nilKeyFormatter) LayerKind() string { return "" }

func (f *nilKeyFormatter) LayerName() string { return f.name }

func (*nilKeyFormatter) FormatKey(_ string) string { return "" }
