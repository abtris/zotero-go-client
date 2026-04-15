package zotero

import (
	"context"
	"encoding/json"
	"io"
	"iter"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Response wraps an HTTP response from the Zotero API.
type Response struct {
	*http.Response

	// LastVersion is the library version from the Last-Modified-Version header.
	LastVersion int

	// TotalResults is the total number of results from the Total-Results header.
	TotalResults int

	// Links holds parsed Link header values.
	Links ResponseLinks
}

// ResponseLinks holds parsed pagination links.
type ResponseLinks struct {
	Next string
	Prev string
	First string
	Last  string
}

// Client manages communication with the Zotero API.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient httpClient
	userAgent  string

	// Services
	Items       *ItemsService
	Collections *CollectionsService
	Searches    *SearchesService
	Tags        *TagsService
	Schema      *SchemaService
	Keys        *KeysService
	Groups      *GroupsService
	FullText    *FullTextService
	Deleted     *DeletedService
}

// NewClient creates a new Zotero API client.
// If apiKey is empty, the ZOTERO_API_KEY environment variable is used.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("ZOTERO_API_KEY")
	}
	base, _ := url.Parse(BaseURL)
	c := &Client{
		baseURL:    base,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		userAgent:  "zotero-go-client/1.0",
	}
	for _, o := range opts {
		o(c)
	}
	c.Items = &ItemsService{client: c}
	c.Collections = &CollectionsService{client: c}
	c.Searches = &SearchesService{client: c}
	c.Tags = &TagsService{client: c}
	c.Schema = &SchemaService{client: c}
	c.Keys = &KeysService{client: c}
	c.Groups = &GroupsService{client: c}
	c.FullText = &FullTextService{client: c}
	c.Deleted = &DeletedService{client: c}
	return c
}

// newRequest creates an API request.
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Zotero-API-Version", APIVersion)
	if c.apiKey != "" {
		req.Header.Set("Zotero-API-Key", c.apiKey)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// do executes an HTTP request and returns the parsed response.
func (c *Client) do(req *http.Request, v any) (*Response, error) {
	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := &Response{Response: httpResp}
	resp.parseHeaders()

	// Handle backoff header
	if backoff := httpResp.Header.Get("Backoff"); backoff != "" {
		if secs, err := strconv.Atoi(backoff); err == nil && secs > 0 {
			time.Sleep(time.Duration(secs) * time.Second)
		}
	}

	if httpResp.StatusCode >= 400 {
		return resp, parseAPIError(httpResp)
	}

	if httpResp.StatusCode == http.StatusNoContent || httpResp.StatusCode == http.StatusNotModified {
		return resp, nil
	}

	if v != nil {
		data, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return resp, err
		}
		if err := json.Unmarshal(data, v); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// parseHeaders extracts Zotero-specific headers from the response.
func (r *Response) parseHeaders() {
	if v := r.Header.Get("Last-Modified-Version"); v != "" {
		r.LastVersion, _ = strconv.Atoi(v)
	}
	if v := r.Header.Get("Total-Results"); v != "" {
		r.TotalResults, _ = strconv.Atoi(v)
	}
	r.Links = parseLinkHeader(r.Header.Get("Link"))
}

// parseLinkHeader parses a Link header into ResponseLinks.
func parseLinkHeader(header string) ResponseLinks {
	var links ResponseLinks
	if header == "" {
		return links
	}
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		segments := strings.SplitN(part, ";", 2)
		if len(segments) != 2 {
			continue
		}
		urlPart := strings.Trim(strings.TrimSpace(segments[0]), "<>")
		relPart := strings.TrimSpace(segments[1])

		if strings.Contains(relPart, `rel="next"`) {
			links.Next = urlPart
		} else if strings.Contains(relPart, `rel="prev"`) {
			links.Prev = urlPart
		} else if strings.Contains(relPart, `rel="first"`) {
			links.First = urlPart
		} else if strings.Contains(relPart, `rel="last"`) {
			links.Last = urlPart
		}
	}
	return links
}

// parseAPIError creates an APIError from an HTTP response.
func parseAPIError(resp *http.Response) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}
	body, err := io.ReadAll(resp.Body)
	if err == nil && len(body) > 0 {
		apiErr.Message = string(body)
	}
	if v := resp.Header.Get("Retry-After"); v != "" {
		apiErr.RetryAfter, _ = strconv.Atoi(v)
	}
	return apiErr
}

// get is a convenience method for GET requests.
func (c *Client) get(ctx context.Context, path string, opts []RequestOption, v any) (*Response, error) {
	params := applyOptions(opts)
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

// post is a convenience method for POST requests.
func (c *Client) post(ctx context.Context, path string, body any, v any) (*Response, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = strings.NewReader(string(data))
	}
	req, err := c.newRequest(ctx, http.MethodPost, path, r)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

// put is a convenience method for PUT requests with a version header.
func (c *Client) put(ctx context.Context, path string, body any, version int, v any) (*Response, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = strings.NewReader(string(data))
	}
	req, err := c.newRequest(ctx, http.MethodPut, path, r)
	if err != nil {
		return nil, err
	}
	if version >= 0 {
		req.Header.Set("If-Unmodified-Since-Version", strconv.Itoa(version))
	}
	return c.do(req, v)
}

// patch is a convenience method for PATCH requests with a version header.
func (c *Client) patch(ctx context.Context, path string, body any, version int, v any) (*Response, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = strings.NewReader(string(data))
	}
	req, err := c.newRequest(ctx, http.MethodPatch, path, r)
	if err != nil {
		return nil, err
	}
	if version >= 0 {
		req.Header.Set("If-Unmodified-Since-Version", strconv.Itoa(version))
	}
	return c.do(req, v)
}

// delete is a convenience method for DELETE requests.
func (c *Client) delete(ctx context.Context, path string, version int) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	if version >= 0 {
		req.Header.Set("If-Unmodified-Since-Version", strconv.Itoa(version))
	}
	return c.do(req, nil)
}

// listAll returns an iterator that fetches all pages of results.
// The fetch function receives the start index and returns the items, response, and error.
func listAll[T any](fetch func(start int) ([]*T, *Response, error)) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		start := 0
		for {
			items, resp, err := fetch(start)
			if err != nil {
				yield(nil, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if resp.Links.Next == "" || len(items) == 0 {
				return
			}
			start += len(items)
		}
	}
}


