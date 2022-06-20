package url

import (
	"net/url"
)

// URL is an unmarshall-able url.
type URL struct {
	url.URL
}
