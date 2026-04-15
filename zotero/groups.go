package zotero

import (
	"context"
	"fmt"
)

// GroupsService handles communication with the group endpoints.
type GroupsService struct {
	client *Client
}

// List returns groups for the given user.
func (s *GroupsService) List(ctx context.Context, userID string, opts ...RequestOption) ([]*Group, *Response, error) {
	var groups []*Group
	path := fmt.Sprintf("/users/%s/groups", userID)
	resp, err := s.client.get(ctx, path, opts, &groups)
	return groups, resp, err
}
