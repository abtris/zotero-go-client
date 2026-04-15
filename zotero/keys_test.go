package zotero

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestKeysCurrent(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/keys/current" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/keys/current")
		}
		info := KeyInfo{
			Key:      "abc123",
			UserID:   12345,
			Username: "testuser",
			Access: KeyAccess{
				User: &KeyPermissions{Library: true, Write: true},
			},
		}
		json.NewEncoder(w).Encode(info)
	})

	ctx := context.Background()
	info, _, err := client.Keys.Current(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if info.Key != "abc123" {
		t.Errorf("Key = %q, want %q", info.Key, "abc123")
	}
	if info.UserID != 12345 {
		t.Errorf("UserID = %d, want 12345", info.UserID)
	}
	if info.Username != "testuser" {
		t.Errorf("Username = %q, want %q", info.Username, "testuser")
	}
	if info.Access.User == nil || !info.Access.User.Write {
		t.Error("expected write access")
	}
}

func TestKeysGet(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/keys/mykey" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/keys/mykey")
		}
		info := KeyInfo{Key: "mykey", UserID: 1}
		json.NewEncoder(w).Encode(info)
	})

	ctx := context.Background()
	info, _, err := client.Keys.Get(ctx, "mykey")
	if err != nil {
		t.Fatal(err)
	}
	if info.Key != "mykey" {
		t.Errorf("Key = %q, want %q", info.Key, "mykey")
	}
}

func TestKeysDelete(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/keys/oldkey" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/keys/oldkey")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Keys.Delete(ctx, "oldkey")
	if err != nil {
		t.Fatal(err)
	}
}
