package url

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	t.Parallel()

	t.Run("create_empty_url", func(t *testing.T) {
		t.Parallel()

		u := URL{}

		assert.Equal(
			t,
			"",
			u.String(),
		)
		assert.Equal(
			t,
			"",
			u.Scheme,
		)
		assert.Equal(
			t,
			"",
			u.Host,
		)
		assert.Equal(
			t,
			"",
			u.Path,
		)
	})

	t.Run("create_url_with_embedded_url", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://example.com/path?query=value")
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		assert.Equal(
			t,
			"https://example.com/path?query=value",
			u.String(),
		)
		assert.Equal(
			t,
			"https",
			u.Scheme,
		)
		assert.Equal(
			t,
			"example.com",
			u.Host,
		)
		assert.Equal(
			t,
			"/path",
			u.Path,
		)
		assert.Equal(
			t,
			"query=value",
			u.RawQuery,
		)
	})

	t.Run("url_methods_work", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse(
			"https://user:pass@example.com:8080/path/to/resource?param=value&other=123#fragment",
		)
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		// Test that all URL methods are accessible.
		assert.Equal(
			t,
			"https",
			u.Scheme,
		)
		assert.Equal(
			t,
			"example.com:8080",
			u.Host,
		)
		assert.Equal(
			t,
			"/path/to/resource",
			u.Path,
		)
		assert.Equal(
			t,
			"param=value&other=123",
			u.RawQuery,
		)
		assert.Equal(
			t,
			"fragment",
			u.Fragment,
		)

		// Test query parsing.
		values := u.Query()
		assert.Equal(
			t,
			"value",
			values.Get("param"),
		)
		assert.Equal(
			t,
			"123",
			values.Get("other"),
		)
	})

	t.Run("different_schemes", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name   string
			urlStr string
			scheme string
		}{
			{
				name:   "http_scheme",
				urlStr: "http://example.com",
				scheme: "http",
			},
			{
				name:   "https_scheme",
				urlStr: "https://secure.example.com",
				scheme: "https",
			},
			{
				name:   "ftp_scheme",
				urlStr: "ftp://files.example.com/file.txt",
				scheme: "ftp",
			},
			{
				name:   "file_scheme",
				urlStr: "file:///path/to/file",
				scheme: "file",
			},
			{
				name:   "custom_scheme",
				urlStr: "myscheme://custom.host/path",
				scheme: "myscheme",
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				baseURL, err := url.Parse(tc.urlStr)
				assert.NoError(t, err)

				u := URL{URL: *baseURL}

				assert.Equal(
					t,
					tc.scheme,
					u.Scheme,
				)
				assert.Equal(
					t,
					tc.urlStr,
					u.String(),
				)
			})
		}
	})

	t.Run("url_modification", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://example.com/path")
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		// Modify the embedded URL.
		u.Path = "/new/path"
		u.RawQuery = "new=query"
		u.Fragment = "newfragment"

		assert.Equal(
			t,
			"/new/path",
			u.Path,
		)
		assert.Equal(
			t,
			"new=query",
			u.RawQuery,
		)
		assert.Equal(
			t,
			"newfragment",
			u.Fragment,
		)
		assert.Equal(
			t,
			"https://example.com/new/path?new=query#newfragment",
			u.String(),
		)
	})

	t.Run("resolve_reference", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://example.com/base/path/")
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		ref, err := url.Parse("../other/resource")
		assert.NoError(t, err)

		resolved := u.ResolveReference(ref)

		assert.Equal(
			t,
			"https://example.com/base/other/resource",
			resolved.String(),
		)
	})

	t.Run("parse_request_uri", func(t *testing.T) {
		t.Parallel()

		uri := "/path/to/resource?param=value"
		baseURL, err := url.ParseRequestURI(uri)
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		assert.Equal(
			t,
			"",
			u.Scheme,
		) // No scheme in request URI.
		assert.Equal(
			t,
			"",
			u.Host,
		) // No host in request URI.
		assert.Equal(
			t,
			"/path/to/resource",
			u.Path,
		)
		assert.Equal(
			t,
			"param=value",
			u.RawQuery,
		)
	})

	t.Run("url_with_port", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://example.com:9000/path")
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		assert.Equal(
			t,
			"example.com:9000",
			u.Host,
		)
		assert.Equal(
			t,
			"9000",
			u.Port(),
		)
		assert.Equal(
			t,
			"example.com",
			u.Hostname(),
		)
	})

	t.Run("url_with_user_info", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse(
			"https://username:password@example.com/secure",
		)
		assert.NoError(t, err)

		u := URL{URL: *baseURL}

		assert.NotNil(t, u.User)
		assert.Equal(
			t,
			"username",
			u.User.Username(),
		)

		password, hasPassword := u.User.Password()
		assert.True(t, hasPassword)
		assert.Equal(
			t,
			"password",
			password,
		)
	})

	t.Run("invalid_url_characters", func(t *testing.T) {
		t.Parallel()

		// Create URL with characters that need escaping.
		baseURL := &url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/path with spaces",
		}

		u := URL{URL: *baseURL}

		assert.Equal(
			t,
			"/path with spaces",
			u.Path,
		)
		assert.Equal(
			t,
			"https://example.com/path%20with%20spaces",
			u.String(),
		)
	})
}
