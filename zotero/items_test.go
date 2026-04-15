package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItemsList(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items")
		}
		w.Header().Set("Last-Modified-Version", "10")
		w.Header().Set("Total-Results", "2")
		items := []Item{
			{Key: "ITEM1", Version: 5, Data: ItemData{ItemType: "book", Title: "Test Book"}},
			{Key: "ITEM2", Version: 6, Data: ItemData{ItemType: "journalArticle", Title: "Test Article"}},
		}
		json.NewEncoder(w).Encode(items)
	})

	ctx := context.Background()
	items, resp, err := client.Items.List(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Key != "ITEM1" {
		t.Errorf("items[0].Key = %q, want %q", items[0].Key, "ITEM1")
	}
	if items[1].Data.Title != "Test Article" {
		t.Errorf("items[1].Data.Title = %q, want %q", items[1].Data.Title, "Test Article")
	}
	if resp.LastVersion != 10 {
		t.Errorf("LastVersion = %d, want 10", resp.LastVersion)
	}
	if resp.TotalResults != 2 {
		t.Errorf("TotalResults = %d, want 2", resp.TotalResults)
	}
}

func TestItemsGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1")
		}
		item := Item{Key: "ITEM1", Version: 5, Data: ItemData{ItemType: "book", Title: "Test Book"}}
		json.NewEncoder(w).Encode(item)
	})

	ctx := context.Background()
	item, _, err := client.Items.Get(ctx, UserLibrary("123"), "ITEM1")
	if err != nil {
		t.Fatal(err)
	}
	if item.Key != "ITEM1" {
		t.Errorf("Key = %q, want %q", item.Key, "ITEM1")
	}
	if item.Data.Title != "Test Book" {
		t.Errorf("Title = %q, want %q", item.Data.Title, "Test Book")
	}
}

func TestItemsCreate(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/users/123/items" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items")
		}
		wr := WriteResponse{
			Success: map[string]string{"0": "NEWITEM1"},
		}
		json.NewEncoder(w).Encode(wr)
	})

	ctx := context.Background()
	items := []*ItemData{{ItemType: "book", Title: "New Book"}}
	wr, _, err := client.Items.Create(ctx, UserLibrary("123"), items)
	if err != nil {
		t.Fatal(err)
	}
	if wr.Success["0"] != "NEWITEM1" {
		t.Errorf("Success[0] = %q, want %q", wr.Success["0"], "NEWITEM1")
	}
}

func TestItemsDelete(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/users/123/items/ITEM1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1")
		}
		if v := r.Header.Get("If-Unmodified-Since-Version"); v != "5" {
			t.Errorf("If-Unmodified-Since-Version = %q, want %q", v, "5")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Items.Delete(ctx, UserLibrary("123"), "ITEM1", 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestItemsListWithOptions(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "title" {
			t.Errorf("sort = %q, want %q", q.Get("sort"), "title")
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit = %q, want %q", q.Get("limit"), "10")
		}
		if q.Get("tag") != "science" {
			t.Errorf("tag = %q, want %q", q.Get("tag"), "science")
		}
		w.Write([]byte("[]"))
	})

	ctx := context.Background()
	_, _, err := client.Items.List(ctx, UserLibrary("123"), WithSort("title"), WithLimit(10), WithTag("science"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestItemsListAll(t *testing.T) {
	callCount := 0
	var ts *httptest.Server
	var client *Client
	client, ts = testServer(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.Header().Set("Link", `<`+ts.URL+`/users/123/items?start=2&limit=100>; rel="next"`)
			items := []Item{
				{Key: "A", Data: ItemData{ItemType: "book"}},
				{Key: "B", Data: ItemData{ItemType: "book"}},
			}
			json.NewEncoder(w).Encode(items)
		} else {
			items := []Item{
				{Key: "C", Data: ItemData{ItemType: "book"}},
			}
			json.NewEncoder(w).Encode(items)
		}
	})

	ctx := context.Background()
	var keys []string
	for item, err := range client.Items.ListAll(ctx, UserLibrary("123")) {
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, item.Key)
	}
	if len(keys) != 3 {
		t.Fatalf("got %d items, want 3", len(keys))
	}
	if keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("keys = %v, want [A B C]", keys)
	}
}
