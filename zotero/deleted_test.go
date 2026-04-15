package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDeletedGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/deleted" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/deleted")
		}
		if q := r.URL.Query().Get("since"); q != "10" {
			t.Errorf("since = %q, want %q", q, "10")
		}
		dc := DeletedContent{
			Collections: []string{"COL1"},
			Searches:    []string{"SRCH1"},
			Items:       []string{"ITEM1", "ITEM2"},
			Tags:        []string{"oldtag"},
			Settings:    []string{},
		}
		json.NewEncoder(w).Encode(dc)
	})

	ctx := context.Background()
	dc, _, err := client.Deleted.Get(ctx, UserLibrary("123"), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(dc.Items) != 2 {
		t.Errorf("got %d deleted items, want 2", len(dc.Items))
	}
	if dc.Items[0] != "ITEM1" {
		t.Errorf("Items[0] = %q, want %q", dc.Items[0], "ITEM1")
	}
	if len(dc.Collections) != 1 || dc.Collections[0] != "COL1" {
		t.Errorf("Collections = %v, want [COL1]", dc.Collections)
	}
	if len(dc.Searches) != 1 || dc.Searches[0] != "SRCH1" {
		t.Errorf("Searches = %v, want [SRCH1]", dc.Searches)
	}
	if len(dc.Tags) != 1 || dc.Tags[0] != "oldtag" {
		t.Errorf("Tags = %v, want [oldtag]", dc.Tags)
	}
}

func TestDeletedGetEmpty(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		dc := DeletedContent{
			Collections: []string{},
			Searches:    []string{},
			Items:       []string{},
			Tags:        []string{},
			Settings:    []string{},
		}
		json.NewEncoder(w).Encode(dc)
	})

	ctx := context.Background()
	dc, _, err := client.Deleted.Get(ctx, UserLibrary("123"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(dc.Items) != 0 {
		t.Errorf("expected no deleted items, got %d", len(dc.Items))
	}
}
