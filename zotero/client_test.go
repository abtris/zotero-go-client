package zotero

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	client := NewClient("test-api-key", WithBaseURL(ts.URL))
	return client, ts
}

func TestNewClient(t *testing.T) {
	c := NewClient("my-key")
	if c.apiKey != "my-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "my-key")
	}
	if c.baseURL.String() != BaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL.String(), BaseURL)
	}
	if c.Items == nil {
		t.Error("Items service is nil")
	}
	if c.Collections == nil {
		t.Error("Collections service is nil")
	}
}

func TestNewClientEnvFallback(t *testing.T) {
	t.Setenv("ZOTERO_API_KEY", "env-key")
	c := NewClient("")
	if c.apiKey != "env-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "env-key")
	}
}

func TestNewClientExplicitKeyOverridesEnv(t *testing.T) {
	t.Setenv("ZOTERO_API_KEY", "env-key")
	c := NewClient("explicit-key")
	if c.apiKey != "explicit-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "explicit-key")
	}
}

func TestRequestHeaders(t *testing.T) {
	var gotHeaders http.Header
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})

	ctx := context.Background()
	client.Items.List(ctx, UserLibrary("12345"))

	if v := gotHeaders.Get("Zotero-API-Version"); v != "3" {
		t.Errorf("Zotero-API-Version = %q, want %q", v, "3")
	}
	if v := gotHeaders.Get("Zotero-API-Key"); v != "test-api-key" {
		t.Errorf("Zotero-API-Key = %q, want %q", v, "test-api-key")
	}
}

func TestResponseHeaders(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified-Version", "42")
		w.Header().Set("Total-Results", "100")
		w.Header().Set("Link", `<https://api.zotero.org/next>; rel="next", <https://api.zotero.org/last>; rel="last"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})

	ctx := context.Background()
	_, resp, err := client.Items.List(ctx, UserLibrary("12345"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.LastVersion != 42 {
		t.Errorf("LastVersion = %d, want 42", resp.LastVersion)
	}
	if resp.TotalResults != 100 {
		t.Errorf("TotalResults = %d, want 100", resp.TotalResults)
	}
	if resp.Links.Next != "https://api.zotero.org/next" {
		t.Errorf("Links.Next = %q, want %q", resp.Links.Next, "https://api.zotero.org/next")
	}
	if resp.Links.Last != "https://api.zotero.org/last" {
		t.Errorf("Links.Last = %q, want %q", resp.Links.Last, "https://api.zotero.org/last")
	}
}

func TestAPIError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Rate limited"))
	})

	ctx := context.Background()
	_, _, err := client.Items.List(ctx, UserLibrary("12345"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsRateLimited(err) {
		t.Errorf("expected rate limited error, got %v", err)
	}
	apiErr := err.(*APIError)
	if apiErr.RetryAfter != 30 {
		t.Errorf("RetryAfter = %d, want 30", apiErr.RetryAfter)
	}
}

func TestNotFoundError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	})

	ctx := context.Background()
	_, _, err := client.Items.Get(ctx, UserLibrary("12345"), "ABCDEFGH")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestParseLinkHeader(t *testing.T) {
	tests := []struct {
		header string
		want   ResponseLinks
	}{
		{
			header: `<https://api.zotero.org/next>; rel="next", <https://api.zotero.org/last>; rel="last"`,
			want:   ResponseLinks{Next: "https://api.zotero.org/next", Last: "https://api.zotero.org/last"},
		},
		{
			header: "",
			want:   ResponseLinks{},
		},
	}
	for _, tt := range tests {
		got := parseLinkHeader(tt.header)
		if got != tt.want {
			t.Errorf("parseLinkHeader(%q) = %+v, want %+v", tt.header, got, tt.want)
		}
	}
}

func TestLibraryID(t *testing.T) {
	u := UserLibrary("12345")
	if u.Path() != "/users/12345" {
		t.Errorf("UserLibrary path = %q, want %q", u.Path(), "/users/12345")
	}
	g := GroupLibrary("67890")
	if g.Path() != "/groups/67890" {
		t.Errorf("GroupLibrary path = %q, want %q", g.Path(), "/groups/67890")
	}
}
