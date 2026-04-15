// Package zotero provides a Go client for the Zotero Web API v3.
//
// Usage:
//
//	client := zotero.NewClient("your-api-key")
//	lib := zotero.UserLibrary("12345")
//	items, resp, err := client.Items.List(ctx, lib)
package zotero

const (
	// BaseURL is the default base URL for the Zotero API.
	BaseURL = "https://api.zotero.org"

	// StreamURL is the WebSocket URL for the Zotero streaming API.
	StreamURL = "wss://stream.zotero.org"

	// APIVersion is the Zotero API version used by this client.
	APIVersion = "3"
)
