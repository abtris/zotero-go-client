package zotero

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// LibraryID identifies a user or group library.
type LibraryID struct {
	prefix string // e.g. "/users/12345" or "/groups/67890"
}

// UserLibrary returns a LibraryID for the given user.
func UserLibrary(userID string) LibraryID {
	return LibraryID{prefix: "/users/" + userID}
}

// GroupLibrary returns a LibraryID for the given group.
func GroupLibrary(groupID string) LibraryID {
	return LibraryID{prefix: "/groups/" + groupID}
}

// Path returns the URL path prefix for this library.
func (l LibraryID) Path() string { return l.prefix }

// RequestOption configures an API list/search request.
type RequestOption func(v url.Values)

// WithLimit sets the maximum number of results to return (1-100).
func WithLimit(n int) RequestOption {
	return func(v url.Values) { v.Set("limit", strconv.Itoa(n)) }
}

// WithStart sets the index of the first result to return.
func WithStart(n int) RequestOption {
	return func(v url.Values) { v.Set("start", strconv.Itoa(n)) }
}

// WithSort sets the field to sort results by.
func WithSort(field string) RequestOption {
	return func(v url.Values) { v.Set("sort", field) }
}

// WithDirection sets the sort direction ("asc" or "desc").
func WithDirection(dir string) RequestOption {
	return func(v url.Values) { v.Set("direction", dir) }
}

// WithFormat sets the response format (e.g., "json", "bib", "csljson").
func WithFormat(format string) RequestOption {
	return func(v url.Values) { v.Set("format", format) }
}

// WithSince returns only objects modified after the given library version.
func WithSince(version int) RequestOption {
	return func(v url.Values) { v.Set("since", strconv.Itoa(version)) }
}

// WithTag filters results by tag. Multiple calls are ANDed.
func WithTag(tag string) RequestOption {
	return func(v url.Values) { v.Add("tag", tag) }
}

// WithQuery sets a quick search query string.
func WithQuery(q string) RequestOption {
	return func(v url.Values) { v.Set("q", q) }
}

// WithQueryMode sets the quick-search mode ("titleCreatorYear" or "everything").
func WithQueryMode(mode string) RequestOption {
	return func(v url.Values) { v.Set("qmode", mode) }
}

// WithItemType filters results by item type (e.g., "book", "-attachment").
func WithItemType(itemType string) RequestOption {
	return func(v url.Values) { v.Set("itemType", itemType) }
}

// WithIncludeTrashed includes trashed items in results.
func WithIncludeTrashed() RequestOption {
	return func(v url.Values) { v.Set("includeTrashed", "1") }
}

// WithItemKey filters by specific item keys (comma-separated or multiple).
func WithItemKey(keys ...string) RequestOption {
	return func(v url.Values) { v.Set("itemKey", strings.Join(keys, ",")) }
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API (useful for testing).
func WithBaseURL(u string) ClientOption {
	return func(c *Client) {
		parsed, err := url.Parse(u)
		if err == nil {
			c.baseURL = parsed
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc httpClient) ClientOption {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) { c.userAgent = ua }
}

// applyOptions applies RequestOptions to a url.Values.
func applyOptions(opts []RequestOption) url.Values {
	v := url.Values{}
	for _, o := range opts {
		o(v)
	}
	return v
}

// httpClient is the interface for making HTTP requests.
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}
