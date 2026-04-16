package zotero

import (
	"context"
	"fmt"
	"iter"
	"strings"
)

// ItemsService handles communication with the item related endpoints.
type ItemsService struct {
	client *Client
}

// List returns items in the given library.
func (s *ItemsService) List(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	resp, err := s.client.get(ctx, lib.Path()+"/items", opts, &items)
	return items, resp, err
}

// ListAll returns an iterator over all items in the library, handling pagination automatically.
func (s *ItemsService) ListAll(ctx context.Context, lib LibraryID, opts ...RequestOption) iter.Seq2[*Item, error] {
	return listAll(func(start int) ([]*Item, *Response, error) {
		o := append(opts, WithStart(start), WithLimit(100))
		return s.List(ctx, lib, o...)
	})
}

// ListTop returns top-level items (no parent) in the library.
func (s *ItemsService) ListTop(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	resp, err := s.client.get(ctx, lib.Path()+"/items/top", opts, &items)
	return items, resp, err
}

// ListTrashed returns trashed items in the library.
func (s *ItemsService) ListTrashed(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	resp, err := s.client.get(ctx, lib.Path()+"/items/trash", opts, &items)
	return items, resp, err
}

// ListChildren returns child items of a given item.
func (s *ItemsService) ListChildren(ctx context.Context, lib LibraryID, itemKey string, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	path := fmt.Sprintf("%s/items/%s/children", lib.Path(), itemKey)
	resp, err := s.client.get(ctx, path, opts, &items)
	return items, resp, err
}

// ListInCollection returns items in a specific collection.
func (s *ItemsService) ListInCollection(ctx context.Context, lib LibraryID, collKey string, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	path := fmt.Sprintf("%s/collections/%s/items", lib.Path(), collKey)
	resp, err := s.client.get(ctx, path, opts, &items)
	return items, resp, err
}

// ListTopInCollection returns top-level items in a specific collection.
func (s *ItemsService) ListTopInCollection(ctx context.Context, lib LibraryID, collKey string, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	path := fmt.Sprintf("%s/collections/%s/items/top", lib.Path(), collKey)
	resp, err := s.client.get(ctx, path, opts, &items)
	return items, resp, err
}

// ListPublications returns items in the user's publications.
func (s *ItemsService) ListPublications(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Item, *Response, error) {
	var items []*Item
	resp, err := s.client.get(ctx, lib.Path()+"/publications/items", opts, &items)
	return items, resp, err
}

// Get returns a single item by key.
func (s *ItemsService) Get(ctx context.Context, lib LibraryID, itemKey string) (*Item, *Response, error) {
	var item Item
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	resp, err := s.client.get(ctx, path, nil, &item)
	if err != nil {
		return nil, resp, err
	}
	return &item, resp, nil
}

// Create creates one or more items in the library.
func (s *ItemsService) Create(ctx context.Context, lib LibraryID, items []*ItemData) (*WriteResponse, *Response, error) {
	var wr WriteResponse
	resp, err := s.client.post(ctx, lib.Path()+"/items", items, &wr)
	if err != nil {
		return nil, resp, err
	}
	return &wr, resp, nil
}

// Update replaces an item (full update).
func (s *ItemsService) Update(ctx context.Context, lib LibraryID, itemKey string, item *ItemData, version int) (*Item, *Response, error) {
	var result Item
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	resp, err := s.client.put(ctx, path, item, version, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// Patch partially updates an item.
func (s *ItemsService) Patch(ctx context.Context, lib LibraryID, itemKey string, fields map[string]any, version int) (*Item, *Response, error) {
	var result Item
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	resp, err := s.client.patch(ctx, path, fields, version, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// Delete deletes a single item by key.
func (s *ItemsService) Delete(ctx context.Context, lib LibraryID, itemKey string, version int) (*Response, error) {
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	return s.client.delete(ctx, path, version)
}

// DeleteMultiple deletes multiple items by key.
func (s *ItemsService) DeleteMultiple(ctx context.Context, lib LibraryID, itemKeys []string, version int) (*Response, error) {
	path := fmt.Sprintf("%s/items?itemKey=%s", lib.Path(), strings.Join(itemKeys, ","))
	return s.client.delete(ctx, path, version)
}

// GetBibliography returns an HTML bibliography for a single item, formatted by the
// Zotero API's built-in citeproc-js engine. Use WithStyle to set the CSL style
// (defaults to "chicago-note-bibliography") and WithLocale for the locale.
func (s *ItemsService) GetBibliography(ctx context.Context, lib LibraryID, itemKey string, opts ...RequestOption) (string, *Response, error) {
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	opts = append([]RequestOption{WithFormat("bib")}, opts...)
	data, resp, err := s.client.getRaw(ctx, path, opts)
	if err != nil {
		return "", resp, err
	}
	return string(data), resp, nil
}

// ListBibliography returns HTML bibliographies for items in the library.
// Each item's bibliography entry is concatenated in the response.
func (s *ItemsService) ListBibliography(ctx context.Context, lib LibraryID, opts ...RequestOption) (string, *Response, error) {
	path := lib.Path() + "/items"
	opts = append([]RequestOption{WithFormat("bib")}, opts...)
	data, resp, err := s.client.getRaw(ctx, path, opts)
	if err != nil {
		return "", resp, err
	}
	return string(data), resp, nil
}

// ListTopBibliography returns HTML bibliographies for top-level items in the library.
func (s *ItemsService) ListTopBibliography(ctx context.Context, lib LibraryID, opts ...RequestOption) (string, *Response, error) {
	path := lib.Path() + "/items/top"
	opts = append([]RequestOption{WithFormat("bib")}, opts...)
	data, resp, err := s.client.getRaw(ctx, path, opts)
	if err != nil {
		return "", resp, err
	}
	return string(data), resp, nil
}

// ListCollectionBibliography returns HTML bibliographies for items in a collection.
func (s *ItemsService) ListCollectionBibliography(ctx context.Context, lib LibraryID, collKey string, opts ...RequestOption) (string, *Response, error) {
	path := fmt.Sprintf("%s/collections/%s/items", lib.Path(), collKey)
	opts = append([]RequestOption{WithFormat("bib")}, opts...)
	data, resp, err := s.client.getRaw(ctx, path, opts)
	if err != nil {
		return "", resp, err
	}
	return string(data), resp, nil
}

// CitationItem holds citation and bibliography data returned when using
// include=citation,bib with format=json.
type CitationItem struct {
	Key      string `json:"key"`
	Version  int    `json:"version"`
	Citation string `json:"citation"`
	Bib      string `json:"bib"`
}

// GetCitation returns the formatted inline citation for a single item.
func (s *ItemsService) GetCitation(ctx context.Context, lib LibraryID, itemKey string, opts ...RequestOption) (string, *Response, error) {
	path := fmt.Sprintf("%s/items/%s", lib.Path(), itemKey)
	opts = append([]RequestOption{
		WithFormat("json"),
		WithInclude("citation"),
	}, opts...)
	var result struct {
		Citation string `json:"citation"`
	}
	resp, err := s.client.get(ctx, path, opts, &result)
	if err != nil {
		return "", resp, err
	}
	return result.Citation, resp, nil
}
