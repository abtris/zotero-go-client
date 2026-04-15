package zotero

import (
	"context"
	"fmt"
	"strings"
)

// SearchesService handles communication with the saved search endpoints.
type SearchesService struct {
	client *Client
}

// List returns saved searches in the given library.
func (s *SearchesService) List(ctx context.Context, lib LibraryID, opts ...RequestOption) ([]*Search, *Response, error) {
	var searches []*Search
	resp, err := s.client.get(ctx, lib.Path()+"/searches", opts, &searches)
	return searches, resp, err
}

// Get returns a single saved search by key.
func (s *SearchesService) Get(ctx context.Context, lib LibraryID, searchKey string) (*Search, *Response, error) {
	var search Search
	path := fmt.Sprintf("%s/searches/%s", lib.Path(), searchKey)
	resp, err := s.client.get(ctx, path, nil, &search)
	if err != nil {
		return nil, resp, err
	}
	return &search, resp, nil
}

// Create creates one or more saved searches.
func (s *SearchesService) Create(ctx context.Context, lib LibraryID, searches []*SearchData) (*WriteResponse, *Response, error) {
	var wr WriteResponse
	resp, err := s.client.post(ctx, lib.Path()+"/searches", searches, &wr)
	if err != nil {
		return nil, resp, err
	}
	return &wr, resp, nil
}

// DeleteMultiple deletes multiple saved searches by key.
func (s *SearchesService) DeleteMultiple(ctx context.Context, lib LibraryID, searchKeys []string, version int) (*Response, error) {
	path := fmt.Sprintf("%s/searches?searchKey=%s", lib.Path(), strings.Join(searchKeys, ","))
	return s.client.delete(ctx, path, version)
}
