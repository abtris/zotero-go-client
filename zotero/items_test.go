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

func TestItemsGetBibliography(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1")
		}
		q := r.URL.Query()
		if q.Get("format") != "bib" {
			t.Errorf("format = %q, want %q", q.Get("format"), "bib")
		}
		if q.Get("style") != "apa" {
			t.Errorf("style = %q, want %q", q.Get("style"), "apa")
		}
		if q.Get("locale") != "en-US" {
			t.Errorf("locale = %q, want %q", q.Get("locale"), "en-US")
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Last-Modified-Version", "10")
		w.Write([]byte(`<div class="csl-bib-body"><div class="csl-entry">Smith, J. (2023). <i>Test Book</i>.</div></div>`))
	})

	ctx := context.Background()
	bib, resp, err := client.Items.GetBibliography(ctx, UserLibrary("123"), "ITEM1", WithStyle("apa"), WithLocale("en-US"))
	if err != nil {
		t.Fatal(err)
	}
	if bib == "" {
		t.Error("bibliography is empty")
	}
	if !contains(bib, "Test Book") {
		t.Errorf("bibliography does not contain 'Test Book': %s", bib)
	}
	if resp.LastVersion != 10 {
		t.Errorf("LastVersion = %d, want 10", resp.LastVersion)
	}
}

func TestItemsListBibliography(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items")
		}
		q := r.URL.Query()
		if q.Get("format") != "bib" {
			t.Errorf("format = %q, want %q", q.Get("format"), "bib")
		}
		w.Write([]byte(`<div class="csl-bib-body"><div class="csl-entry">Entry 1</div><div class="csl-entry">Entry 2</div></div>`))
	})

	ctx := context.Background()
	bib, _, err := client.Items.ListBibliography(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(bib, "Entry 1") || !contains(bib, "Entry 2") {
		t.Errorf("bibliography missing entries: %s", bib)
	}
}

func TestItemsListTopBibliography(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/top" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/top")
		}
		q := r.URL.Query()
		if q.Get("format") != "bib" {
			t.Errorf("format = %q, want %q", q.Get("format"), "bib")
		}
		w.Write([]byte(`<div class="csl-bib-body"><div class="csl-entry">Top Entry</div></div>`))
	})

	ctx := context.Background()
	bib, _, err := client.Items.ListTopBibliography(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(bib, "Top Entry") {
		t.Errorf("bibliography missing 'Top Entry': %s", bib)
	}
}

func TestItemsListCollectionBibliography(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections/COL1/items" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/collections/COL1/items")
		}
		q := r.URL.Query()
		if q.Get("format") != "bib" {
			t.Errorf("format = %q, want %q", q.Get("format"), "bib")
		}
		w.Write([]byte(`<div class="csl-bib-body"><div class="csl-entry">Collection Entry</div></div>`))
	})

	ctx := context.Background()
	bib, _, err := client.Items.ListCollectionBibliography(ctx, UserLibrary("123"), "COL1")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(bib, "Collection Entry") {
		t.Errorf("bibliography missing 'Collection Entry': %s", bib)
	}
}

func TestItemsGetCitation(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/items/ITEM1")
		}
		q := r.URL.Query()
		if q.Get("format") != "json" {
			t.Errorf("format = %q, want %q", q.Get("format"), "json")
		}
		if q.Get("include") != "citation" {
			t.Errorf("include = %q, want %q", q.Get("include"), "citation")
		}
		if q.Get("style") != "apa" {
			t.Errorf("style = %q, want %q", q.Get("style"), "apa")
		}
		w.Write([]byte(`{"citation":"(Smith, 2023)"}`))
	})

	ctx := context.Background()
	citation, _, err := client.Items.GetCitation(ctx, UserLibrary("123"), "ITEM1", WithStyle("apa"))
	if err != nil {
		t.Fatal(err)
	}
	if citation != "(Smith, 2023)" {
		t.Errorf("citation = %q, want %q", citation, "(Smith, 2023)")
	}
}

func TestItemsGetBibliographyWithLinkWrap(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("linkwrap") != "1" {
			t.Errorf("linkwrap = %q, want %q", q.Get("linkwrap"), "1")
		}
		if q.Get("format") != "bib" {
			t.Errorf("format = %q, want %q", q.Get("format"), "bib")
		}
		w.Write([]byte(`<div class="csl-bib-body"><div class="csl-entry"><a href="https://example.com">Link</a></div></div>`))
	})

	ctx := context.Background()
	bib, _, err := client.Items.GetBibliography(ctx, UserLibrary("123"), "ITEM1", WithLinkWrap())
	if err != nil {
		t.Fatal(err)
	}
	if !contains(bib, "https://example.com") {
		t.Errorf("bibliography missing link: %s", bib)
	}
}

func TestItemsGetBibliographyError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	})

	ctx := context.Background()
	_, _, err := client.Items.GetBibliography(ctx, UserLibrary("123"), "BADKEY")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

// contains is a test helper to check substring presence.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
