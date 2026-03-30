/*
Package model provides reflection-based struct analysis and field mapping for
dsco's configuration processing system.

# Overview

The model package implements the core reflection logic that analyzes Go structs
and creates an internal representation (model) for configuration processing.
This model is used to map configuration values from various sources (command
line, environment variables, files) to struct fields, handle type conversion,
and validate configuration completeness.

# Core Type

## Model

The Model struct represents an analyzed configuration structure:

	type Model struct {
		accelerator Node              // Tree structure for field traversal
		getList     GetListInterface  // Operations for value retrieval
		expandList  ExpandListInterface // Operations for struct expansion
		typeName    string            // Full type name for error messages
		fieldCount  uint              // Total number of configurable fields
	}

# Model Creation

## NewModel Function

NewModel analyzes a reflect.Type and builds the internal model:

	func NewModel(inputModelType reflect.Type) (*Model, error)

### Analysis Process

1. **Struct Scanning**: Recursively scans the struct type and all nested types
2. **Field Mapping**: Creates unique identifiers (UIDs) for each configurable field
3. **Tag Processing**: Extracts YAML tags for key name mapping
4. **Visibility Rules**: Applies Go visibility rules for embedded fields
5. **Error Collection**: Aggregates all model building errors

### Example Usage

	type Config struct {
		Host     *string         `yaml:"host"`
		Port     *int            `yaml:"port"`
		Database *DatabaseConfig `yaml:"database"`
	}

	model, err := model.NewModel(reflect.TypeOf(&Config{}))
	if err != nil {
		log.Fatal(err)
	}

# Field Processing

## Struct Field Analysis

The package analyzes struct fields for configuration mapping:

### Supported Field Types

- **Pointer types**: *string, *int, *bool, *time.Duration, etc.
- **Slice types**: []string, []int, etc.
- **Map types**: map[string]string, map[string]interface{}, etc.
- **Nested structs**: *DatabaseConfig, *ServerConfig, etc.
- **Embedded structs**: Anonymous struct fields

### Field Visibility Rules

Go's visibility rules are applied:

	type Config struct {
		PublicField  *string `yaml:"public"`   // Included
		privateField *string `yaml:"private"`  // Excluded (unexported)
	}

### Embedded Struct Handling

Embedded structs are flattened following Go's visibility rules:

	type Base struct {
		Host *string `yaml:"host"`
	}

	type Config struct {
		Base                           // Embedded: Host accessible as Config.Host
		Port *int `yaml:"port"`
	}

# Value Operations

## ApplyOn Method

Applies a value getter to extract configuration values:

	func (m *Model) ApplyOn(g ValueGetter) (fvalue.Values, error)

This method is used by string-based layers (cmdline, env) to extract
values that match the model's field paths.

## Expand Method

Expands struct-based configuration sources:

	func (m *Model) Expand(g StructExpander) error

Used for struct layers to traverse and extract values from Go structs.

## Fill Method

Populates the target struct with resolved values:

	func (m *Model) Fill(
		inputModelValue reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error)

This is the final step where resolved configuration values are assigned
to struct fields.

## GetFieldValuesFor Method

Extracts field values from a struct instance:

	func (m *Model) GetFieldValuesFor(
		id string,
		value reflect.Value,
	) fvalue.Values

Used by struct layers to create fvalue.Values from Go structs.

# Node System

## Node Interface

The model uses a tree structure of nodes for efficient field access:

	type Node interface {
		BuildGetList(list *GetList)
		BuildExpandList(list *ExpandList)
		FeedFieldValues(id string, values fvalue.Values, value reflect.Value)
		Fill(value reflect.Value, layers []fvalue.Values) (plocation.Locations, error)
	}

### Node Types

- **StructNode**: Represents struct types with nested fields
- **ValueNode**: Represents leaf fields (actual configuration values)
- **SliceNode**: Represents slice/array fields
- **MapNode**: Represents map fields

# Field Path Generation

## Path Format

Configuration paths use dot notation for nested fields:

	host                    # Top-level field
	database.host           # Nested field
	database.pool.size      # Deeply nested field
	servers.primary.port    # Complex nesting

## YAML Tag Processing

YAML tags determine configuration keys:

	type Config struct {
		DatabaseHost *string `yaml:"db_host"`      // Key: "db_host"
		MaxRetries   *int    `yaml:"max-retries"`  // Key: "max-retries"
		Timeout      *int                          // Key: "timeout" (from field name)
	}

# Error Handling

## ModelError Type

Model building errors are aggregated:

	type ModelError struct {
		merror.MError
	}

	var ErrModel = errors.New("")

### Error Types

- **InvalidEmbedded**: Embedded pointer structs not supported
- **FieldNameCollision**: Multiple fields resolve to same name
- **UnsupportedType**: Field type not supported for configuration

### Error Checking

	model, err := model.NewModel(configType)
	if err != nil {
		var modelErr model.ModelError
		if errors.As(err, &modelErr) {
			for _, e := range modelErr.Errors() {
				log.Printf("Model error: %v", e)
			}
		}
	}

# Embedded Struct Processing

## Visibility Rules

Embedded fields follow Go's promotion rules:

	type Base struct {
		Host *string `yaml:"host"`
		Port *int    `yaml:"port"`
	}

	type Extended struct {
		Host *string `yaml:"host"`  // Shadows Base.Host
	}

	type Config struct {
		Base                         // Host, Port promoted
		Extended                     // Host shadows Base.Host
		Timeout *int `yaml:"timeout"`
	}

	// Resulting fields:
	// - host (from Extended, shadows Base)
	// - port (from Base)
	// - timeout (from Config)

## Field Name Collision Detection

Collisions at the same depth level generate errors:

	type A struct {
		Field *string `yaml:"field"`
	}

	type B struct {
		Field *string `yaml:"field"`
	}

	type Config struct {
		A  // Both have "field" at same depth
		B  // Error: FieldNameCollisionError
	}

# Integration with dsco Components

## Layer Processing

Models are used by the filler to process configuration layers:

	func (c *dscoContext) generateModel() {
		model, err := model.NewModel(reflect.TypeOf(c.inputModelRef).Elem())
		if err != nil {
			c.err.Add(err)
			return
		}
		c.model = model
	}

## Value Extraction

Each layer uses the model for value extraction:

	func (l *envLayer) GetFieldValuesFrom(model ModelInterface) (fvalue.Values, error) {
		return model.ApplyOn(l.valueGetter)
	}

## Struct Filling

The fill operation uses the model's node tree:

	locations, err := c.model.Fill(targetValue, c.layerFieldValues)

# Performance Considerations

## Model Caching

Models are computed once per configuration type:

- Reflection operations occur during model creation
- Field traversal uses pre-computed node tree
- UID-based lookups are O(1) during filling

## Memory Usage

Model memory is proportional to struct complexity:

- One node per struct field
- Path strings stored once per field
- UID maps for fast field lookup

## Struct Scanning

Scanning occurs depth-first:

- Stack-based traversal avoids recursion limits
- Embedded struct fields processed in declaration order
- Early termination on critical errors

# Thread Safety

Model instances are safe for concurrent use after creation:

- Immutable structure after NewModel returns
- No shared mutable state during Fill operations
- Safe for concurrent configuration processing

# Best Practices

## Configuration Struct Design

Design structs for clear configuration mapping:

	// Good: Clear hierarchy and naming
	type Config struct {
		Server   *ServerConfig   `yaml:"server"`
		Database *DatabaseConfig `yaml:"database"`
		Logging  *LoggingConfig  `yaml:"logging"`
	}

	// Less clear: Flat structure with prefixed names
	type Config struct {
		ServerHost     *string `yaml:"server_host"`
		ServerPort     *int    `yaml:"server_port"`
		DatabaseHost   *string `yaml:"database_host"`
		DatabasePort   *int    `yaml:"database_port"`
	}

## YAML Tag Conventions

Use consistent YAML tag naming:

	// Good: Consistent snake_case or kebab-case
	type Config struct {
		MaxRetries    *int `yaml:"max_retries"`
		ConnTimeout   *int `yaml:"conn_timeout"`
		ReadTimeout   *int `yaml:"read_timeout"`
	}

	// Inconsistent: Mixed naming styles
	type Config struct {
		MaxRetries    *int `yaml:"maxRetries"`
		ConnTimeout   *int `yaml:"conn-timeout"`
		ReadTimeout   *int `yaml:"read_timeout"`
	}

## Embedded Struct Usage

Use embedded structs for shared configuration:

	type CommonConfig struct {
		Timeout *time.Duration `yaml:"timeout"`
		Retries *int           `yaml:"retries"`
	}

	type ServiceConfig struct {
		CommonConfig                    // Shared fields
		Endpoint *string `yaml:"endpoint"`
	}

The model package provides the foundational reflection and analysis
capabilities that enable dsco's type-safe configuration processing
across all supported configuration sources.
*/
package model
