package zotero

import "encoding/json"

// Item represents a Zotero library item (book, article, note, etc.).
type Item struct {
	Key     string          `json:"key"`
	Version int             `json:"version"`
	Library Library         `json:"library"`
	Links   json.RawMessage `json:"links,omitempty"`
	Meta    json.RawMessage `json:"meta,omitempty"`
	Data    ItemData        `json:"data"`
}

// ItemData holds the mutable fields of an item.
// Fields vary by ItemType; use the map-based accessor for dynamic fields.
type ItemData struct {
	Key          string            `json:"key,omitempty"`
	Version      int               `json:"version,omitempty"`
	ItemType     string            `json:"itemType"`
	Title        string            `json:"title,omitempty"`
	Creators     []Creator         `json:"creators,omitempty"`
	AbstractNote string            `json:"abstractNote,omitempty"`
	Date         string            `json:"date,omitempty"`
	URL          string            `json:"url,omitempty"`
	Tags         []Tag             `json:"tags,omitempty"`
	Collections  []string          `json:"collections,omitempty"`
	Relations    map[string]any    `json:"relations,omitempty"`
	ParentItem   string            `json:"parentItem,omitempty"`
	Note         string            `json:"note,omitempty"`
	Extra        map[string]string `json:"-"` // overflow fields
}

// Creator represents an author or other creator of an item.
type Creator struct {
	CreatorType string `json:"creatorType"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Name        string `json:"name,omitempty"` // single-field mode
}

// Tag represents a tag attached to an item.
type Tag struct {
	Tag  string `json:"tag"`
	Type int    `json:"type,omitempty"`
}

// TagEntry represents a tag returned from tag listing endpoints.
type TagEntry struct {
	Tag   string          `json:"tag"`
	Links json.RawMessage `json:"links,omitempty"`
	Meta  *TagMeta        `json:"meta,omitempty"`
}

// TagMeta holds metadata for a tag entry.
type TagMeta struct {
	Type     int `json:"type"`
	NumItems int `json:"numItems"`
}

// Collection represents a Zotero collection.
type Collection struct {
	Key     string          `json:"key"`
	Version int             `json:"version"`
	Library Library         `json:"library"`
	Links   json.RawMessage `json:"links,omitempty"`
	Meta    json.RawMessage `json:"meta,omitempty"`
	Data    CollectionData  `json:"data"`
}

// CollectionData holds the mutable fields of a collection.
type CollectionData struct {
	Key              string         `json:"key,omitempty"`
	Version          int            `json:"version,omitempty"`
	Name             string         `json:"name"`
	ParentCollection any            `json:"parentCollection,omitempty"` // string key or false
	Relations        map[string]any `json:"relations,omitempty"`
}

// Search represents a saved search.
type Search struct {
	Key     string          `json:"key"`
	Version int             `json:"version"`
	Library Library         `json:"library"`
	Links   json.RawMessage `json:"links,omitempty"`
	Data    SearchData      `json:"data"`
}

// SearchData holds the mutable fields of a saved search.
type SearchData struct {
	Key        string            `json:"key,omitempty"`
	Version    int               `json:"version,omitempty"`
	Name       string            `json:"name"`
	Conditions []SearchCondition `json:"conditions"`
}

// SearchCondition is a single condition within a saved search.
type SearchCondition struct {
	Condition string `json:"condition"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

// Library identifies the library an object belongs to.
type Library struct {
	Type  string          `json:"type"`
	ID    int             `json:"id"`
	Name  string          `json:"name"`
	Links json.RawMessage `json:"links,omitempty"`
}

// Group represents a Zotero group.
type Group struct {
	ID      int             `json:"id"`
	Version int             `json:"version"`
	Links   json.RawMessage `json:"links,omitempty"`
	Meta    json.RawMessage `json:"meta,omitempty"`
	Data    GroupData       `json:"data"`
}

// GroupData holds the fields of a group.
type GroupData struct {
	ID                int    `json:"id"`
	Version           int    `json:"version"`
	Name              string `json:"name"`
	Owner             int    `json:"owner"`
	Type              string `json:"type"`
	Description       string `json:"description"`
	URL               string `json:"url"`
	LibraryEditing    string `json:"libraryEditing"`
	LibraryReading    string `json:"libraryReading"`
	FileEditing       string `json:"fileEditing"`
}

// FullText represents the full-text content of an item.
type FullText struct {
	Content    string `json:"content,omitempty"`
	IndexedPages int  `json:"indexedPages,omitempty"`
	TotalPages   int  `json:"totalPages,omitempty"`
	IndexedChars int  `json:"indexedChars,omitempty"`
	Version      int  `json:"version,omitempty"`
}

// FullTextVersions maps item keys to their full-text content versions.
type FullTextVersions map[string]int

// DeletedContent represents deleted library content since a given version.
type DeletedContent struct {
	Collections []string `json:"collections"`
	Searches    []string `json:"searches"`
	Items       []string `json:"items"`
	Tags        []string `json:"tags"`
	Settings    []string `json:"settings"`
}

// KeyInfo represents API key information.
type KeyInfo struct {
	Key      string         `json:"key"`
	UserID   int            `json:"userID"`
	Username string         `json:"username"`
	Access   KeyAccess      `json:"access"`
}

// KeyAccess describes the permissions granted by an API key.
type KeyAccess struct {
	User   *KeyPermissions            `json:"user,omitempty"`
	Groups map[string]*KeyPermissions `json:"groups,omitempty"`
}

// KeyPermissions describes permissions for a library.
type KeyPermissions struct {
	Library bool `json:"library"`
	Files   bool `json:"files,omitempty"`
	Notes   bool `json:"notes,omitempty"`
	Write   bool `json:"write,omitempty"`
}

// SchemaItemType represents an item type from the schema endpoint.
type SchemaItemType struct {
	ItemType  string `json:"itemType"`
	Localized string `json:"localized"`
}

// SchemaField represents a field from the schema endpoint.
type SchemaField struct {
	Field     string `json:"field"`
	Localized string `json:"localized"`
}

// SchemaCreatorType represents a creator type from the schema endpoint.
type SchemaCreatorType struct {
	CreatorType string `json:"creatorType"`
	Localized   string `json:"localized"`
}
