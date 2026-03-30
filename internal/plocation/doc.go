/*
Package plocation provides path location tracking for dsco's configuration
debugging and reporting system.

# Overview

The plocation package defines types for tracking where configuration values
originated during the filling process. This information is essential for
debugging configuration issues, understanding precedence resolution, and
providing clear error messages that help users identify the source of
configuration problems.

# Core Types

## Location

Location represents a single configuration value's source information:

	type Location struct {
		Path     string  // Configuration path (e.g., "database.host")
		Location string  // Source location (e.g., "env[MYAPP-DATABASE-HOST]")
		UID      uint    // Unique identifier for the field
	}

### Fields

- **Path**: The configuration path identifying the field within the struct
  hierarchy. Uses dot notation for nested fields (e.g., "server.timeout").

- **Location**: A human-readable string indicating where the value came from.
  Format varies by source type (cmdline, env, struct, file).

- **UID**: A unique identifier assigned during model building that corresponds
  to a specific field in the configuration structure.

## Locations

Locations is a slice of Location values representing all resolved fields:

	type Locations []Location

# Location Formats

Different configuration sources use consistent location formats:

## Command Line

	"cmdline[--host]"
	"cmdline[--database-port]"
	"cmdline[--max-retry-count]"

## Environment Variables

	"env[MYAPP-HOST]"
	"env[MYAPP-DATABASE-HOST]"
	"env[PREFIX-NESTED-FIELD]"

## Struct Sources

	"struct[defaults]:host"
	"struct[production]:database.timeout"
	"struct[base]:server.port"

## File Sources

	"file[config.yaml]:host"
	"kfile[/etc/config]:DATABASE-HOST"
	"file[/app/settings.json]:api.key"

# Usage Examples

## Retrieving Locations

The Fill function returns locations for all resolved values:

	locations, err := dsco.Fill(
		&config,
		dsco.WithStructLayer(defaults, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
		dsco.WithCmdlineLayer(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Inspect where values came from
	for _, loc := range locations {
		fmt.Printf("%s = (from %s)\n", loc.Path, loc.Location)
	}

## Debugging Configuration

Use locations to understand configuration resolution:

	// Output example:
	// host = (from cmdline[--host])
	// port = (from env[MYAPP-PORT])
	// database.host = (from struct[defaults]:database.host)
	// database.port = (from env[MYAPP-DATABASE-PORT])
	// timeout = (from struct[defaults]:timeout)

## Dump Method

Write a formatted report to any io.Writer:

	func (f *Locations) Dump(writer io.Writer)

### Example Output

	locations.Dump(os.Stdout)

	// Output:
	//   path            |  Location
	//   ----            |  --------
	//   host            |  cmdline[--host]
	//   port            |  env[MYAPP-PORT]
	//   database.host   |  struct[defaults]:database.host
	//   timeout         |  struct[defaults]:timeout

# Location Management

## Report Method

Add a new location entry during filling:

	func (f *Locations) Report(uid uint, path string, location string)

Used internally by the filling process:

	locations.Report(fieldUID, "database.host", "env[MYAPP-DATABASE-HOST]")

## Append Method

Combine location sets:

	func (f *Locations) Append(other Locations)

Used when processing nested structures:

	parentLocations.Append(nestedLocations)

# Integration with dsco Components

## Fill Process

Locations are populated during the Fill operation:

	func (m *Model) Fill(
		inputModelValue reflect.Value,
		layers []fvalue.Values,
	) (plocation.Locations, error) {
		var locations plocation.Locations

		// For each field, track its source
		for uid, field := range m.fields {
			value, source := resolveValue(uid, layers)
			if value != nil {
				locations.Report(uid, field.Path, source.Location)
				// ... assign value ...
			}
		}

		return locations, nil
	}

## Error Reporting

Location information enhances error messages:

	func reportOverrideError(locations plocation.Locations, overriddenUID uint) error {
		for _, loc := range locations {
			if loc.UID == overriddenUID {
				return fmt.Errorf(
					"field '%s' from %s was overridden",
					loc.Path,
					loc.Location,
				)
			}
		}
		return nil
	}

## Configuration Audit

Create audit logs of configuration sources:

	func auditConfiguration(locations plocation.Locations) {
		log.Println("Configuration sources:")
		for _, loc := range locations {
			log.Printf("  %s: %s", loc.Path, loc.Location)
		}
	}

# Dump Output Format

## Tabular Format

The Dump method uses tabwriter for aligned output:

	locations.Dump(os.Stdout)

	// Produces:
	//   path              |  Location
	//   ----              |  --------
	//   host              |  cmdline[--host]
	//   port              |  env[MYAPP-PORT]
	//   database.host     |  struct[defaults]:database.host

## Custom Output

For custom formatting, iterate over locations:

	for _, loc := range locations {
		fmt.Printf("%-20s <- %s\n", loc.Path, loc.Location)
	}

	// Produces:
	// host                 <- cmdline[--host]
	// port                 <- env[MYAPP-PORT]

# Path Conventions

## Simple Fields

Top-level fields use their yaml tag or field name:

	host        # From yaml:"host" or field Host
	port        # From yaml:"port" or field Port
	timeout     # From yaml:"timeout" or field Timeout

## Nested Fields

Nested struct fields use dot notation:

	database.host      # Database struct, Host field
	server.timeout     # Server struct, Timeout field
	api.auth.token     # Deeply nested field

## Array/Slice Fields

Array elements may include indices (future enhancement):

	servers[0].host    # First server's host
	endpoints[1].url   # Second endpoint's URL

# Best Practices

## Logging Configuration Sources

Log configuration sources at startup:

	func logConfiguration(locations plocation.Locations) {
		log.Println("Application configuration:")
		for _, loc := range locations {
			// Don't log sensitive values, just sources
			log.Printf("  %s: loaded from %s", loc.Path, loc.Location)
		}
	}

## Configuration Documentation

Generate configuration documentation:

	func documentConfiguration(locations plocation.Locations) string {
		var sb strings.Builder
		sb.WriteString("# Configuration Sources\n\n")

		for _, loc := range locations {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", loc.Path, loc.Location))
		}

		return sb.String()
	}

## Debugging Overrides

Identify unexpected configuration overrides:

	func checkForOverrides(expected, actual plocation.Locations) {
		expectedSources := make(map[string]string)
		for _, loc := range expected {
			expectedSources[loc.Path] = loc.Location
		}

		for _, loc := range actual {
			if expected, ok := expectedSources[loc.Path]; ok {
				if expected != loc.Location {
					log.Printf("WARNING: %s expected from %s, got from %s",
						loc.Path, expected, loc.Location)
				}
			}
		}
	}

# Thread Safety

Locations types have specific thread safety characteristics:

- **Location struct**: Immutable after creation, safe for concurrent reads
- **Locations slice**: Not safe for concurrent modification
- **Report/Append methods**: Not safe for concurrent use

The filling process operates sequentially, so thread safety is managed
by the calling code.

# Memory Considerations

## Slice Growth

Locations slice grows during filling:

- Initial capacity based on model field count
- Efficient append operations
- Memory proportional to number of configured fields

## String Storage

Location strings are stored directly:

- Path strings typically short (field paths)
- Location strings include source information
- No deduplication of repeated location strings

# Performance Characteristics

## Report Operation

Adding locations is O(1) amortized:

	locations.Report(uid, path, location)

## Append Operation

Appending is O(n) where n is the appended slice length:

	locations.Append(otherLocations)

## Dump Operation

Dump is O(n) where n is number of locations:

	locations.Dump(writer)

The plocation package provides essential debugging and audit capabilities
for dsco's configuration system, enabling users to understand exactly where
their configuration values originated and troubleshoot any issues effectively.
*/
package plocation
