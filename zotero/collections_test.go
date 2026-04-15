package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCollectionsList(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/collections")
		}
		colls := []Collection{
			{Key: "COL1", Version: 3, Data: CollectionData{Name: "My Collection"}},
		}
		json.NewEncoder(w).Encode(colls)
	})

	ctx := context.Background()
	colls, _, err := client.Collections.List(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if len(colls) != 1 {
		t.Fatalf("got %d collections, want 1", len(colls))
	}
	if colls[0].Data.Name != "My Collection" {
		t.Errorf("Name = %q, want %q", colls[0].Data.Name, "My Collection")
	}
}

func TestCollectionsGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections/COL1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/collections/COL1")
		}
		coll := Collection{Key: "COL1", Version: 3, Data: CollectionData{Name: "My Collection"}}
		json.NewEncoder(w).Encode(coll)
	})

	ctx := context.Background()
	coll, _, err := client.Collections.Get(ctx, UserLibrary("123"), "COL1")
	if err != nil {
		t.Fatal(err)
	}
	if coll.Key != "COL1" {
		t.Errorf("Key = %q, want %q", coll.Key, "COL1")
	}
}

func TestCollectionsCreate(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		wr := WriteResponse{
			Success: map[string]string{"0": "NEWCOL1"},
		}
		json.NewEncoder(w).Encode(wr)
	})

	ctx := context.Background()
	colls := []*CollectionData{{Name: "New Collection"}}
	wr, _, err := client.Collections.Create(ctx, UserLibrary("123"), colls)
	if err != nil {
		t.Fatal(err)
	}
	if wr.Success["0"] != "NEWCOL1" {
		t.Errorf("Success[0] = %q, want %q", wr.Success["0"], "NEWCOL1")
	}
}

func TestCollectionsDelete(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if v := r.Header.Get("If-Unmodified-Since-Version"); v != "3" {
			t.Errorf("If-Unmodified-Since-Version = %q, want %q", v, "3")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Collections.Delete(ctx, UserLibrary("123"), "COL1", 3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCollectionsSubcollections(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections/COL1/collections" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/collections/COL1/collections")
		}
		colls := []Collection{
			{Key: "SUB1", Data: CollectionData{Name: "Subcollection"}},
		}
		json.NewEncoder(w).Encode(colls)
	})

	ctx := context.Background()
	colls, _, err := client.Collections.GetSubcollections(ctx, UserLibrary("123"), "COL1")
	if err != nil {
		t.Fatal(err)
	}
	if len(colls) != 1 {
		t.Fatalf("got %d, want 1", len(colls))
	}
	if colls[0].Key != "SUB1" {
		t.Errorf("Key = %q, want %q", colls[0].Key, "SUB1")
	}
}
