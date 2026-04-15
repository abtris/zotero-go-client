package zotero

import (
	"context"
	"encoding/json"
	"net/url"
)

// SchemaService handles communication with the Zotero schema endpoints.
type SchemaService struct {
	client *Client
}

// ItemTypes returns the list of available item types.
func (s *SchemaService) ItemTypes(ctx context.Context) ([]*SchemaItemType, *Response, error) {
	var types []*SchemaItemType
	resp, err := s.client.get(ctx, "/itemTypes", nil, &types)
	return types, resp, err
}

// ItemFields returns the list of all item fields.
func (s *SchemaService) ItemFields(ctx context.Context) ([]*SchemaField, *Response, error) {
	var fields []*SchemaField
	resp, err := s.client.get(ctx, "/itemFields", nil, &fields)
	return fields, resp, err
}

// ItemTypeFields returns the valid fields for a given item type.
func (s *SchemaService) ItemTypeFields(ctx context.Context, itemType string) ([]*SchemaField, *Response, error) {
	var fields []*SchemaField
	resp, err := s.client.get(ctx, "/itemTypeFields", []RequestOption{
		func(v url.Values) { v.Set("itemType", itemType) },
	}, &fields)
	return fields, resp, err
}

// ItemTypeCreatorTypes returns the valid creator types for a given item type.
func (s *SchemaService) ItemTypeCreatorTypes(ctx context.Context, itemType string) ([]*SchemaCreatorType, *Response, error) {
	var types []*SchemaCreatorType
	resp, err := s.client.get(ctx, "/itemTypeCreatorTypes", []RequestOption{
		func(v url.Values) { v.Set("itemType", itemType) },
	}, &types)
	return types, resp, err
}

// CreatorFields returns the list of creator fields.
func (s *SchemaService) CreatorFields(ctx context.Context) ([]*SchemaField, *Response, error) {
	var fields []*SchemaField
	resp, err := s.client.get(ctx, "/creatorFields", nil, &fields)
	return fields, resp, err
}

// ItemTemplate returns a template for creating a new item of the given type.
func (s *SchemaService) ItemTemplate(ctx context.Context, itemType string) (json.RawMessage, *Response, error) {
	var tmpl json.RawMessage
	resp, err := s.client.get(ctx, "/items/new", []RequestOption{
		func(v url.Values) { v.Set("itemType", itemType) },
	}, &tmpl)
	return tmpl, resp, err
}
