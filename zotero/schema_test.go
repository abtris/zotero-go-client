package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSchemaItemTypes(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/itemTypes" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/itemTypes")
		}
		types := []SchemaItemType{
			{ItemType: "book", Localized: "Book"},
			{ItemType: "journalArticle", Localized: "Journal Article"},
		}
		json.NewEncoder(w).Encode(types)
	})

	ctx := context.Background()
	types, _, err := client.Schema.ItemTypes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 {
		t.Fatalf("got %d types, want 2", len(types))
	}
	if types[0].ItemType != "book" {
		t.Errorf("ItemType = %q, want %q", types[0].ItemType, "book")
	}
}

func TestSchemaItemFields(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/itemFields" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/itemFields")
		}
		fields := []SchemaField{
			{Field: "title", Localized: "Title"},
		}
		json.NewEncoder(w).Encode(fields)
	})

	ctx := context.Background()
	fields, _, err := client.Schema.ItemFields(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Field != "title" {
		t.Errorf("unexpected fields: %+v", fields)
	}
}

func TestSchemaItemTypeFields(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/itemTypeFields" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/itemTypeFields")
		}
		if q := r.URL.Query().Get("itemType"); q != "book" {
			t.Errorf("itemType = %q, want %q", q, "book")
		}
		fields := []SchemaField{{Field: "title", Localized: "Title"}}
		json.NewEncoder(w).Encode(fields)
	})

	ctx := context.Background()
	fields, _, err := client.Schema.ItemTypeFields(ctx, "book")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Fatalf("got %d fields, want 1", len(fields))
	}
}

func TestSchemaItemTypeCreatorTypes(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/itemTypeCreatorTypes" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/itemTypeCreatorTypes")
		}
		if q := r.URL.Query().Get("itemType"); q != "book" {
			t.Errorf("itemType = %q, want %q", q, "book")
		}
		types := []SchemaCreatorType{{CreatorType: "author", Localized: "Author"}}
		json.NewEncoder(w).Encode(types)
	})

	ctx := context.Background()
	types, _, err := client.Schema.ItemTypeCreatorTypes(ctx, "book")
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 1 || types[0].CreatorType != "author" {
		t.Errorf("unexpected creator types: %+v", types)
	}
}

func TestSchemaItemTemplate(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/items/new" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/items/new")
		}
		if q := r.URL.Query().Get("itemType"); q != "book" {
			t.Errorf("itemType = %q, want %q", q, "book")
		}
		w.Write([]byte(`{"itemType":"book","title":"","creators":[]}`))
	})

	ctx := context.Background()
	tmpl, _, err := client.Schema.ItemTemplate(ctx, "book")
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]any
	if err := json.Unmarshal(tmpl, &data); err != nil {
		t.Fatal(err)
	}
	if data["itemType"] != "book" {
		t.Errorf("itemType = %v, want %q", data["itemType"], "book")
	}
}
