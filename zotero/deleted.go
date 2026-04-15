package zotero

import "context"

// DeletedService handles communication with the deleted content endpoint.
type DeletedService struct {
	client *Client
}

// Get returns content deleted from the library since the given version.
func (s *DeletedService) Get(ctx context.Context, lib LibraryID, since int) (*DeletedContent, *Response, error) {
	var dc DeletedContent
	resp, err := s.client.get(ctx, lib.Path()+"/deleted", []RequestOption{
		WithSince(since),
	}, &dc)
	if err != nil {
		return nil, resp, err
	}
	return &dc, resp, nil
}
