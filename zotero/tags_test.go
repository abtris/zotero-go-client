package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTagsList(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/tags" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/tags")
		}
		tags := []TagEntry{
			{Tag: "science", Meta: &TagMeta{Type: 0, NumItems: 5}},
			{Tag: "history", Meta: &TagMeta{Type: 0, NumItems: 3}},
		}
		json.NewEncoder(w).Encode(tags)
	})

	ctx := context.Background()
	tags, _, err := client.Tags.List(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("got %d tags, want 2", len(tags))
	}
	if tags[0].Tag != "science" {
		t.Errorf("Tag = %q, want %q", tags[0].Tag, "science")
	}
	if tags[0].Meta.NumItems != 5 {
		t.Errorf("NumItems = %d, want 5", tags[0].Meta.NumItems)
	}
}

func TestTagsGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/tags/science" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/tags/science")
		}
		tags := []TagEntry{{Tag: "science"}}
		json.NewEncoder(w).Encode(tags)
	})

	ctx := context.Background()
	tags, _, err := client.Tags.Get(ctx, UserLibrary("123"), "science")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags[0].Tag != "science" {
		t.Errorf("unexpected tags: %+v", tags)
	}
}

func TestTagsListForItem(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1/tags" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1/tags")
		}
		tags := []TagEntry{{Tag: "tagged"}}
		json.NewEncoder(w).Encode(tags)
	})

	ctx := context.Background()
	tags, _, err := client.Tags.ListForItem(ctx, UserLibrary("123"), "ITEM1")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 {
		t.Fatalf("got %d tags, want 1", len(tags))
	}
}

func TestTagsListInCollection(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections/COL1/tags" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/collections/COL1/tags")
		}
		json.NewEncoder(w).Encode([]TagEntry{})
	})

	ctx := context.Background()
	_, _, err := client.Tags.ListInCollection(ctx, UserLibrary("123"), "COL1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTagsDeleteMultiple(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if v := r.Header.Get("If-Unmodified-Since-Version"); v != "5" {
			t.Errorf("If-Unmodified-Since-Version = %q, want %q", v, "5")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Tags.DeleteMultiple(ctx, UserLibrary("123"), []string{"tag1", "tag2"}, 5)
	if err != nil {
		t.Fatal(err)
	}
}
