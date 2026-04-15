package zotero

import (
	"context"
	"fmt"
	"iter"
	"strings"
)

// CollectionsService handles communication with the collection related endpoints.
type CollectionsService struct {
	client *Client
}

// List returns collections in the given library.
func (s *CollectionsService) List(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Collection, *Response, error) {
	var colls []*Collection
	resp, err := s.client.get(ctx, lib.Path()+"/collections", opts, &colls)
	return colls, resp, err
}

// ListAll returns an iterator over all collections, handling pagination automatically.
func (s *CollectionsService) ListAll(ctx context.Context, lib LibraryID, opts ...RequestOption) iter.Seq2[*Collection, error] {
	return listAll(func(start int) ([]*Collection, *Response, error) {
		o := append(opts, WithStart(start), WithLimit(100))
		return s.List(ctx, lib, o...)
	})
}

// ListTop returns top-level collections (no parent) in the library.
func (s *CollectionsService) ListTop(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Collection, *Response, error) {
	var colls []*Collection
	resp, err := s.client.get(ctx, lib.Path()+"/collections/top", opts, &colls)
	return colls, resp, err
}

// Get returns a single collection by key.
func (s *CollectionsService) Get(ctx context.Context, lib LibraryID, collKey string) (*Collection, *Response, error) {
	var coll Collection
	path := fmt.Sprintf("%s/collections/%s", lib.Path(), collKey)
	resp, err := s.client.get(ctx, path, nil, &coll)
	if err != nil {
		return nil, resp, err
	}
	return &coll, resp, nil
}

// GetSubcollections returns subcollections of a given collection.
func (s *CollectionsService) GetSubcollections(ctx context.Context, lib LibraryID, collKey string, opts ...RequestOption) ([]*Collection, *Response, error) {
	var colls []*Collection
	path := fmt.Sprintf("%s/collections/%s/collections", lib.Path(), collKey)
	resp, err := s.client.get(ctx, path, opts, &colls)
	return colls, resp, err
}

// Create creates one or more collections in the library.
func (s *CollectionsService) Create(ctx context.Context, lib LibraryID, colls []*CollectionData) (*WriteResponse, *Response, error) {
	var wr WriteResponse
	resp, err := s.client.post(ctx, lib.Path()+"/collections", colls, &wr)
	if err != nil {
		return nil, resp, err
	}
	return &wr, resp, nil
}

// Update replaces a collection (full update).
func (s *CollectionsService) Update(ctx context.Context, lib LibraryID, collKey string, coll *CollectionData, version int) (*Collection, *Response, error) {
	var result Collection
	path := fmt.Sprintf("%s/collections/%s", lib.Path(), collKey)
	resp, err := s.client.put(ctx, path, coll, version, &result)
	if err != nil {
		return nil, resp, err
	}
	return &result, resp, nil
}

// Delete deletes a single collection by key.
func (s *CollectionsService) Delete(ctx context.Context, lib LibraryID, collKey string, version int) (*Response, error) {
	path := fmt.Sprintf("%s/collections/%s", lib.Path(), collKey)
	return s.client.delete(ctx, path, version)
}

// DeleteMultiple deletes multiple collections by key.
func (s *CollectionsService) DeleteMultiple(ctx context.Context, lib LibraryID, collKeys []string, version int) (*Response, error) {
	path := fmt.Sprintf("%s/collections?collectionKey=%s", lib.Path(), strings.Join(collKeys, ","))
	return s.client.delete(ctx, path, version)
}
