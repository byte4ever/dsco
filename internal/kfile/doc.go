/*
Package kfile provides file-based configuration value extraction for dsco's
configuration system.

# Overview

The kfile package implements a file-based configuration provider that reads
configuration values from files within a directory structure. Each file in
the directory represents a single configuration key, with the file name
serving as the key and the file contents as the value. This approach is
commonly used in Kubernetes ConfigMaps and Secrets mounted as volumes.

# File Naming Convention

Files must follow a strict naming pattern to be recognized:

	Pattern: ^[A-Z][A-Z\d]*([-_][A-Z][A-Z\d]*)*$

## Valid File Names

	HOST                    # Simple key
	DATABASE-HOST           # Kebab-case key
	API_KEY                 # Snake_case key
	MAX-RETRY-COUNT         # Multiple segments
	DB_CONNECTION_POOL      # Mixed separators

## Invalid File Names

	host                    # Lowercase not allowed
	123-KEY                 # Cannot start with digit
	-INVALID                # Cannot start with separator
	KEY--DOUBLE             # No consecutive separators

# Directory Structure

The package recursively processes files in the specified directory:

	/config/
	├── HOST                # → key: "host"
	├── PORT                # → key: "port"
	├── DATABASE-HOST       # → key: "database-host"
	├── DATABASE-PORT       # → key: "database-port"
	└── API_KEY             # → key: "api_key"

# Key Transformation

File names are transformed to lowercase configuration keys:

	HOST              → "host"
	DATABASE-HOST     → "database-host"
	MAX_RETRY_COUNT   → "max_retry_count"
	API-KEY-V2        → "api-key-v2"

# Usage Examples

## Basic Usage

	provider, err := kfile.NewEntriesProvider("/etc/myapp/config")
	if err != nil {
		log.Fatal(err)
	}

	values := provider.GetStringValues()
	for key, value := range values {
		fmt.Printf("%s = %s (from %s)\n", key, value.Value, value.Location)
	}

## Integration with dsco

	import "github.com/byte4ever/dsco"

	type Config struct {
		Host     *string `yaml:"host"`
		Port     *int    `yaml:"port"`
		APIKey   *string `yaml:"api_key"`
	}

	func main() {
		var config *Config

		fileProvider, err := kfile.NewEntriesProvider("/etc/myapp/config")
		if err != nil {
			log.Fatal(err)
		}

		_, err = dsco.Fill(
			&config,
			dsco.WithStringValueProvider(fileProvider),
			dsco.WithEnvLayer("MYAPP"),
			dsco.WithCmdlineLayer(),
		)

		if err != nil {
			log.Fatal(err)
		}
	}

## Kubernetes Usage

Ideal for Kubernetes ConfigMaps and Secrets:

	# ConfigMap mounted at /etc/config
	apiVersion: v1
	kind: ConfigMap
	metadata:
		name: myapp-config
	data:
		HOST: "api.production.com"
		PORT: "8080"
		LOG-LEVEL: "info"

	# Secret mounted at /etc/secrets
	apiVersion: v1
	kind: Secret
	metadata:
		name: myapp-secrets
	data:
		API-KEY: <base64-encoded>
		DATABASE-PASSWORD: <base64-encoded>

	# Pod spec
	volumes:
		- name: config
		  configMap:
			name: myapp-config
		- name: secrets
		  secret:
			secretName: myapp-secrets

Go code:

	// Load from both config and secrets directories
	configProvider, _ := kfile.NewEntriesProvider("/etc/config")
	secretsProvider, _ := kfile.NewEntriesProvider("/etc/secrets")

	_, err := dsco.Fill(
		&config,
		dsco.WithStringValueProvider(configProvider),
		dsco.WithStringValueProvider(secretsProvider),
	)

# Location Tracking

Each value includes its file path for debugging:

	Location format: "kfile[directory]:filename"

Examples:

	"kfile[/etc/config]:HOST"
	"kfile[/etc/secrets]:API-KEY"
	"kfile[/app/config]:DATABASE-HOST"

# Error Handling

## Invalid File Name Errors

Files with invalid names generate errors:

	var pathErrs kfile.PathErrors
	if errors.As(err, &pathErrs) {
		for _, pathErr := range pathErrs {
			fmt.Printf("Invalid file: %s - %v\n", pathErr.Path, pathErr.Err)
		}
	}

## File Read Errors

File access errors are reported with full context:

	"kfile[/etc/config]:SECRET - permission denied"
	"kfile[/etc/config]:DATA - file too large"

# File System Abstraction

The package uses afero for file system abstraction, enabling:

## Production Usage

Uses the real OS file system:

	provider, err := kfile.NewEntriesProvider("/etc/config")

## Testing

Can use in-memory file systems for testing (internal):

	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/config/HOST", []byte("localhost"), 0644)
	// Internal: provider, err := newProvider(fs, "/config", opts)

# File Content Handling

## Content Reading

File contents are read entirely as string values:

	# File: /config/MESSAGE
	Hello, World!

	# Results in:
	# key: "message"
	# value: "Hello, World!"

## Whitespace Preservation

Leading and trailing whitespace in file contents is preserved:

	# File: /config/TEMPLATE
	  indented content

	# value: "  indented content\n"

## Binary Safety

Files are read as bytes and converted to string. Binary content
may not be handled correctly.

# Provider Interface

## EntriesProvider Type

The provider implements dsco's NamedStringValuesProvider interface:

	type EntriesProvider struct {
		values svalue.Values
		name   string
	}

	func (e *EntriesProvider) GetName() string
	func (e *EntriesProvider) GetStringValues() svalue.Values

## GetName Method

Returns a descriptive name for error messages:

	provider.GetName()  // "kfile(/etc/config)"

## GetStringValues Method

Returns all configuration values from the directory:

	values := provider.GetStringValues()
	// Returns svalue.Values map

# Options and Configuration

## Silent File Errors

Internal option to suppress file read errors (useful for optional files):

	opts := &options{
		silentFileErrors: true,
	}

# Best Practices

## Directory Organization

Organize configuration files logically:

	/etc/myapp/
	├── config/           # General configuration
	│   ├── HOST
	│   ├── PORT
	│   └── LOG-LEVEL
	├── database/         # Database configuration
	│   ├── HOST
	│   ├── PORT
	│   └── NAME
	└── secrets/          # Sensitive configuration
		├── API-KEY
		└── DATABASE-PASSWORD

## Naming Conventions

Use consistent naming across environments:

	# Development
	/dev/config/DATABASE-HOST → "localhost"

	# Production
	/prod/config/DATABASE-HOST → "db.production.internal"

## Error Handling

Always check for errors when loading file configurations:

	provider, err := kfile.NewEntriesProvider(configPath)
	if err != nil {
		var pathErrs kfile.PathErrors
		if errors.As(err, &pathErrs) {
			// Handle specific file errors
			for _, e := range pathErrs {
				log.Printf("Config file error: %s: %v", e.Path, e.Err)
			}
		}
		return err
	}

# Thread Safety

The package is designed for safe concurrent use:

- EntriesProvider instances are immutable after creation
- GetStringValues() is safe for concurrent calls
- Multiple providers can be used concurrently
- File reading occurs only during provider creation

# Performance Considerations

## Startup Time

File reading occurs once during provider creation:

- All files are read during NewEntriesProvider call
- Values are cached in memory
- No file I/O during configuration resolution

## Memory Usage

Memory proportional to total file content size:

- Each file content stored as string
- Location strings add minor overhead
- Empty directories create minimal overhead

## Large Files

Avoid using kfile for large configuration files:

- Entire file content loaded into memory
- Consider using YAML/JSON file providers for structured data
- kfile is optimal for small, single-value configurations

The kfile package provides seamless integration with file-based configuration
systems like Kubernetes ConfigMaps and Secrets, enabling secure and flexible
configuration management in containerized environments.
*/
package kfile
