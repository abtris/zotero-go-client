package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGroupsList(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/groups" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/groups")
		}
		groups := []Group{
			{ID: 100, Version: 1, Data: GroupData{
				ID:   100,
				Name: "Test Group",
				Type: "PublicOpen",
			}},
			{ID: 200, Version: 2, Data: GroupData{
				ID:   200,
				Name: "Private Group",
				Type: "Private",
			}},
		}
		json.NewEncoder(w).Encode(groups)
	})

	ctx := context.Background()
	groups, _, err := client.Groups.List(ctx, "123")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}
	if groups[0].Data.Name != "Test Group" {
		t.Errorf("Name = %q, want %q", groups[0].Data.Name, "Test Group")
	}
	if groups[1].Data.Type != "Private" {
		t.Errorf("Type = %q, want %q", groups[1].Data.Type, "Private")
	}
}

func TestGroupsListEmpty(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[]"))
	})

	ctx := context.Background()
	groups, _, err := client.Groups.List(ctx, "999")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Errorf("got %d groups, want 0", len(groups))
	}
}

func TestGroupsListGroupLibrary(t *testing.T) {
	// Verify GroupLibrary helper creates correct path prefix for other services
	lib := GroupLibrary("100")
	if lib.Path() != "/groups/100" {
		t.Errorf("GroupLibrary path = %q, want %q", lib.Path(), "/groups/100")
	}

	// Use group library to list items
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups/100/items" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/groups/100/items")
		}
		w.Write([]byte("[]"))
	})

	ctx := context.Background()
	items, _, err := client.Items.List(ctx, GroupLibrary("100"))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Errorf("got %d items, want 0", len(items))
	}
}
