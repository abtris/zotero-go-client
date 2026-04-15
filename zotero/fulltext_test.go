package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestFullTextListChanged(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/fulltext" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/fulltext")
		}
		if q := r.URL.Query().Get("since"); q != "5" {
			t.Errorf("since = %q, want %q", q, "5")
		}
		versions := FullTextVersions{"ITEM1": 10, "ITEM2": 12}
		json.NewEncoder(w).Encode(versions)
	})

	ctx := context.Background()
	versions, _, err := client.FullText.ListChanged(ctx, UserLibrary("123"), 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 2 {
		t.Fatalf("got %d versions, want 2", len(versions))
	}
	if versions["ITEM1"] != 10 {
		t.Errorf("ITEM1 version = %d, want 10", versions["ITEM1"])
	}
}

func TestFullTextGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1/fulltext" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1/fulltext")
		}
		ft := FullText{
			Content:      "This is the full text content.",
			IndexedPages: 5,
			TotalPages:   10,
		}
		json.NewEncoder(w).Encode(ft)
	})

	ctx := context.Background()
	ft, _, err := client.FullText.Get(ctx, UserLibrary("123"), "ITEM1")
	if err != nil {
		t.Fatal(err)
	}
	if ft.Content != "This is the full text content." {
		t.Errorf("Content = %q, want %q", ft.Content, "This is the full text content.")
	}
	if ft.IndexedPages != 5 {
		t.Errorf("IndexedPages = %d, want 5", ft.IndexedPages)
	}
}

func TestFullTextSet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/users/123/items/ITEM1/fulltext" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1/fulltext")
		}
		var ft FullText
		if err := json.NewDecoder(r.Body).Decode(&ft); err != nil {
			t.Fatal(err)
		}
		if ft.Content != "new content" {
			t.Errorf("Content = %q, want %q", ft.Content, "new content")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.FullText.Set(ctx, UserLibrary("123"), "ITEM1", &FullText{
		Content:      "new content",
		IndexedPages: 1,
		TotalPages:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
}
