package zotero

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// TagsService handles communication with the tag related endpoints.
type TagsService struct {
	client *Client
}

// List returns tags in the given library.
func (s *TagsService) List(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	resp, err := s.client.get(ctx, lib.Path()+"/tags", opts, &tags)
	return tags, resp, err
}

// Get returns a single tag by name.
func (s *TagsService) Get(ctx context.Context, lib LibraryID, tag string) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	path := fmt.Sprintf("%s/tags/%s", lib.Path(), url.PathEscape(tag))
	resp, err := s.client.get(ctx, path, nil, &tags)
	return tags, resp, err
}

// ListForItem returns tags for a specific item.
func (s *TagsService) ListForItem(ctx context.Context, lib LibraryID, itemKey string, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	path := fmt.Sprintf("%s/items/%s/tags", lib.Path(), itemKey)
	resp, err := s.client.get(ctx, path, opts, &tags)
	return tags, resp, err
}

// ListInCollection returns tags used in a specific collection.
func (s *TagsService) ListInCollection(ctx context.Context, lib LibraryID, collKey string, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	path := fmt.Sprintf("%s/collections/%s/tags", lib.Path(), collKey)
	resp, err := s.client.get(ctx, path, opts, &tags)
	return tags, resp, err
}

// ListForItems returns all tags used across items in the library.
func (s *TagsService) ListForItems(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	resp, err := s.client.get(ctx, lib.Path()+"/items/tags", opts, &tags)
	return tags, resp, err
}

// ListForTopItems returns tags used by top-level items.
func (s *TagsService) ListForTopItems(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	resp, err := s.client.get(ctx, lib.Path()+"/items/top/tags", opts, &tags)
	return tags, resp, err
}

// ListForTrashedItems returns tags used by trashed items.
func (s *TagsService) ListForTrashedItems(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*TagEntry, *Response, error) {
	var tags []*TagEntry
	resp, err := s.client.get(ctx, lib.Path()+"/items/trash/tags", opts, &tags)
	return tags, resp, err
}

// DeleteMultiple deletes multiple tags by name.
func (s *TagsService) DeleteMultiple(ctx context.Context, lib LibraryID, tags []string, version int) (*Response, error) {
	encoded := make([]string, len(tags))
	for i, t := range tags {
		encoded[i] = url.QueryEscape(t)
	}
	path := fmt.Sprintf("%s/tags?tag=%s", lib.Path(), strings.Join(encoded, "+||+"))
	return s.client.delete(ctx, path, version)
}
