/*
Package url provides an unmarshallable URL type for dsco configuration fields.

# Overview

The url package defines a URL type that wraps Go's standard net/url.URL with
YAML unmarshaling capabilities. This allows URLs to be used directly as
configuration field types in dsco structs, with automatic parsing from string
configuration values.

# Core Type

## URL

URL embeds net/url.URL and inherits all its methods:

	type URL struct {
		url.URL
	}

The URL type can be used anywhere a standard url.URL is expected, while also
supporting direct unmarshaling from configuration sources.

# Usage in Configuration

## Basic Usage

Use URL as a pointer field in configuration structs:

	type Config struct {
		APIEndpoint *url.URL `yaml:"api_endpoint"`
		WebhookURL  *url.URL `yaml:"webhook_url"`
		ProxyURL    *url.URL `yaml:"proxy_url"`
	}

	var config *Config
	_, err := dsco.Fill(
		&config,
		dsco.WithStructLayer(&Config{
			APIEndpoint: parseURL("https://api.example.com/v1"),
			ProxyURL:    parseURL("http://proxy.internal:8080"),
		}, "defaults"),
		dsco.WithEnvLayer("MYAPP"),
	)

## Configuration Sources

URLs can be provided from various configuration sources:

### Environment Variables

	MYAPP-API-ENDPOINT=https://api.production.com/v2
	MYAPP-WEBHOOK-URL=https://hooks.example.com/callback
	MYAPP-PROXY-URL=http://10.0.0.1:3128

### Command Line

	./myapp --api-endpoint=https://api.staging.com --proxy-url=socks5://localhost:1080

### Struct Defaults

	defaults := &Config{
		APIEndpoint: &url.URL{URL: url.URL{
			Scheme: "https",
			Host:   "api.example.com",
			Path:   "/v1",
		}},
	}

# URL Parsing

The package handles URL parsing following standard URL semantics:

## Supported URL Schemes

	http://example.com/path
	https://secure.example.com:8443/api
	ftp://files.example.com/downloads
	file:///local/path/to/file
	postgres://user:pass@host:5432/database
	redis://cache.internal:6379/0
	amqp://rabbitmq.svc:5672/vhost

## URL Components

All standard URL components are supported:

	scheme://[user:password@]host[:port]/path[?query][#fragment]

Examples:

	https://api.example.com                    # Simple HTTPS
	https://api.example.com:8443               # With port
	https://api.example.com/v1/users           # With path
	https://api.example.com/search?q=test      # With query
	https://user:pass@api.example.com          # With credentials
	https://api.example.com/path#section       # With fragment

# Integration with dsco

## Type Conversion

URLs are automatically converted from string configuration values:

	// Configuration source provides: "https://api.example.com/v1"
	// dsco converts to: &url.URL{
	//     Scheme: "https",
	//     Host:   "api.example.com",
	//     Path:   "/v1",
	// }

## Error Handling

Invalid URLs produce clear error messages:

	// Invalid URL: "not a valid url"
	// Error: "field 'api_endpoint': invalid URL format"

## Location Tracking

URL values include source location for debugging:

	// Location: "env[MYAPP-API-ENDPOINT]"
	// Location: "cmdline[--api-endpoint]"

# URL Operations

## Accessing URL Components

Once configured, URL fields can be used like standard url.URL:

	config := loadConfig()

	// Get host and port
	host := config.APIEndpoint.Host
	port := config.APIEndpoint.Port()

	// Get full URL string
	urlString := config.APIEndpoint.String()

	// Build request
	resp, err := http.Get(config.APIEndpoint.String() + "/users")

## URL Manipulation

Create derived URLs from configured base:

	baseURL := config.APIEndpoint

	// Build path-specific URLs
	usersURL := *baseURL
	usersURL.Path = path.Join(baseURL.Path, "users")

	ordersURL := *baseURL
	ordersURL.Path = path.Join(baseURL.Path, "orders")

	// Add query parameters
	searchURL := *baseURL
	searchURL.RawQuery = "q=search+term&limit=10"

# Best Practices

## URL Validation

Validate URLs beyond simple parsing when needed:

	func validateAPIEndpoint(u *url.URL) error {
		if u == nil {
			return errors.New("API endpoint is required")
		}

		if u.Scheme != "https" {
			return errors.New("API endpoint must use HTTPS")
		}

		if u.Host == "" {
			return errors.New("API endpoint must have a host")
		}

		return nil
	}

## Default URL Construction

Use helper functions for default URL creation:

	func mustParseURL(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			panic(fmt.Sprintf("invalid URL: %s", rawURL))
		}
		return &url.URL{URL: *u}
	}

	defaults := &Config{
		APIEndpoint: mustParseURL("https://api.example.com/v1"),
	}

## Sensitive URLs

Handle URLs with credentials carefully:

	// URLs may contain credentials
	// postgres://user:password@host:5432/db

	// Log URLs without credentials
	func sanitizeURL(u *url.URL) string {
		sanitized := *u
		sanitized.User = nil
		return sanitized.String()
	}

# Thread Safety

URL instances are safe for concurrent read access but not for concurrent
modification. In typical dsco usage, URLs are configured once at startup
and then read-only, making them safe for concurrent use.

# Error Scenarios

## Common URL Errors

The following URL issues are detected and reported:

	// Missing scheme
	"example.com/path"  // Error: missing scheme

	// Invalid characters
	"https://example.com/path with spaces"  // May need encoding

	// Invalid port
	"https://example.com:invalid/path"  // Error: invalid port

## Error Messages

Errors include context for debugging:

	"field 'api_endpoint' from env[MYAPP-API-ENDPOINT]: invalid URL format"
	"field 'proxy_url' from cmdline[--proxy-url]: missing scheme"

# Testing

When testing with URL configurations:

	func TestConfig(t *testing.T) {
		config := &Config{
			APIEndpoint: mustParseURL("https://test.example.com"),
		}

		_, err := dsco.Fill(
			&config,
			dsco.WithStructLayer(config, "test"),
		)

		require.NoError(t, err)
		assert.Equal(t, "test.example.com", config.APIEndpoint.Host)
	}

The url package provides seamless integration of URL configuration in dsco,
enabling type-safe URL handling with automatic parsing and validation.
*/
package url
