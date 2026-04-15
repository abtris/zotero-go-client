package zotero

import (
	"context"
	"fmt"
)

// FullTextService handles communication with the full-text content endpoints.
type FullTextService struct {
	client *Client
}

// ListChanged returns item keys with full-text content changed since the given version.
func (s *FullTextService) ListChanged(ctx context.Context, lib LibraryID, since int) (FullTextVersions, *Response, error) {
	var versions FullTextVersions
	resp, err := s.client.get(ctx, lib.Path()+"/fulltext", []RequestOption{
		WithSince(since),
	}, &versions)
	return versions, resp, err
}

// Get returns the full-text content for an item.
func (s *FullTextService) Get(ctx context.Context, lib LibraryID, itemKey string) (*FullText, *Response, error) {
	var ft FullText
	path := fmt.Sprintf("%s/items/%s/fulltext", lib.Path(), itemKey)
	resp, err := s.client.get(ctx, path, nil, &ft)
	if err != nil {
		return nil, resp, err
	}
	return &ft, resp, nil
}

// Set sets the full-text content for an item.
func (s *FullTextService) Set(ctx context.Context, lib LibraryID, itemKey string, ft *FullText) (*Response, error) {
	path := fmt.Sprintf("%s/items/%s/fulltext", lib.Path(), itemKey)
	return s.client.put(ctx, path, ft, -1, nil)
}
