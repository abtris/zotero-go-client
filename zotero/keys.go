package zotero

import (
	"context"
	"fmt"
)

// KeysService handles communication with the API key endpoints.
type KeysService struct {
	client *Client
}

// Current returns information about the current API key.
func (s *KeysService) Current(ctx context.Context) (*KeyInfo, *Response, error) {
	var info KeyInfo
	resp, err := s.client.get(ctx, "/keys/current", nil, &info)
	if err != nil {
		return nil, resp, err
	}
	return &info, resp, nil
}

// Get returns information about a specific API key.
func (s *KeysService) Get(ctx context.Context, key string) (*KeyInfo, *Response, error) {
	var info KeyInfo
	path := fmt.Sprintf("/keys/%s", key)
	resp, err := s.client.get(ctx, path, nil, &info)
	if err != nil {
		return nil, resp, err
	}
	return &info, resp, nil
}

// Delete deletes a specific API key.
func (s *KeysService) Delete(ctx context.Context, key string) (*Response, error) {
	path := fmt.Sprintf("/keys/%s", key)
	return s.client.delete(ctx, path, -1)
}
