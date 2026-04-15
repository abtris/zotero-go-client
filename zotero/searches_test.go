package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSearchesList(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/searches" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/searches")
		}
		searches := []Search{
			{Key: "SRCH1", Version: 1, Data: SearchData{Name: "My Search", Conditions: []SearchCondition{
				{Condition: "title", Operator: "contains", Value: "test"},
			}}},
		}
		json.NewEncoder(w).Encode(searches)
	})

	ctx := context.Background()
	searches, _, err := client.Searches.List(ctx, UserLibrary("123"))
	if err != nil {
		t.Fatal(err)
	}
	if len(searches) != 1 {
		t.Fatalf("got %d searches, want 1", len(searches))
	}
	if searches[0].Data.Name != "My Search" {
		t.Errorf("Name = %q, want %q", searches[0].Data.Name, "My Search")
	}
	if len(searches[0].Data.Conditions) != 1 {
		t.Fatalf("got %d conditions, want 1", len(searches[0].Data.Conditions))
	}
	if searches[0].Data.Conditions[0].Value != "test" {
		t.Errorf("Condition.Value = %q, want %q", searches[0].Data.Conditions[0].Value, "test")
	}
}

func TestSearchesGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/searches/SRCH1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/users/123/searches/SRCH1")
		}
		search := Search{Key: "SRCH1", Data: SearchData{Name: "My Search"}}
		json.NewEncoder(w).Encode(search)
	})

	ctx := context.Background()
	search, _, err := client.Searches.Get(ctx, UserLibrary("123"), "SRCH1")
	if err != nil {
		t.Fatal(err)
	}
	if search.Key != "SRCH1" {
		t.Errorf("Key = %q, want %q", search.Key, "SRCH1")
	}
}

func TestSearchesCreate(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body []SearchData
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body) != 1 || body[0].Name != "New Search" {
			t.Errorf("unexpected body: %+v", body)
		}
		wr := WriteResponse{Success: map[string]string{"0": "NEWSRCH"}}
		json.NewEncoder(w).Encode(wr)
	})

	ctx := context.Background()
	data := []*SearchData{{
		Name: "New Search",
		Conditions: []SearchCondition{{Condition: "title", Operator: "is", Value: "foo"}},
	}}
	wr, _, err := client.Searches.Create(ctx, UserLibrary("123"), data)
	if err != nil {
		t.Fatal(err)
	}
	if wr.Success["0"] != "NEWSRCH" {
		t.Errorf("Success[0] = %q, want %q", wr.Success["0"], "NEWSRCH")
	}
}

func TestSearchesDeleteMultiple(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if v := r.URL.Query().Get("searchKey"); v != "S1,S2" {
			t.Errorf("searchKey = %q, want %q", v, "S1,S2")
		}
		if v := r.Header.Get("If-Unmodified-Since-Version"); v != "10" {
			t.Errorf("If-Unmodified-Since-Version = %q, want %q", v, "10")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Searches.DeleteMultiple(ctx, UserLibrary("123"), []string{"S1", "S2"}, 10)
	if err != nil {
		t.Fatal(err)
	}
}
