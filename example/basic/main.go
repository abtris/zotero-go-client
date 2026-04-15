// Example program demonstrating the Zotero Go client library.
//
// Usage:
//
//	export ZOTERO_API_KEY="your-api-key"
//	export ZOTERO_USER_ID="your-user-id"
//	go run ./example/basic
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abtris/zotero-go-client/zotero"
)

func main() {
	userID := os.Getenv("ZOTERO_USER_ID")
	if userID == "" {
		log.Fatal("Set ZOTERO_USER_ID environment variable")
	}

	// NewClient("") automatically reads ZOTERO_API_KEY from the environment.
	client := zotero.NewClient("")
	lib := zotero.UserLibrary(userID)
	ctx := context.Background()

	// --- Verify API key ---
	keyInfo, _, err := client.Keys.Current(ctx)
	if err != nil {
		log.Fatalf("Keys.Current: %v", err)
	}
	fmt.Printf("Authenticated as: %s (user %d)\n\n", keyInfo.Username, keyInfo.UserID)

	// --- List collections ---
	fmt.Println("=== Collections ===")
	colls, _, err := client.Collections.List(ctx, lib, zotero.WithLimit(5))
	if err != nil {
		log.Fatalf("Collections.List: %v", err)
	}
	for _, c := range colls {
		fmt.Printf("  [%s] %s\n", c.Key, c.Data.Name)
	}
	fmt.Println()

	// --- List top-level items ---
	fmt.Println("=== Top Items (first 10) ===")
	items, resp, err := client.Items.ListTop(ctx, lib, zotero.WithLimit(10), zotero.WithSort("dateModified"), zotero.WithDirection("desc"))
	if err != nil {
		log.Fatalf("Items.ListTop: %v", err)
	}
	fmt.Printf("Total items in library: %d\n", resp.TotalResults)
	for _, item := range items {
		creators := ""
		if len(item.Data.Creators) > 0 {
			c := item.Data.Creators[0]
			if c.Name != "" {
				creators = c.Name
			} else {
				creators = c.LastName + ", " + c.FirstName
			}
		}
		fmt.Printf("  [%s] %-15s %s — %s\n", item.Key, item.Data.ItemType, item.Data.Title, creators)
	}
	fmt.Println()

	// --- Iterate all items (pagination demo) ---
	fmt.Println("=== Counting all items via iterator ===")
	count := 0
	for _, err := range client.Items.ListAll(ctx, lib) {
		if err != nil {
			log.Fatalf("Items.ListAll: %v", err)
		}
		count++
	}
	fmt.Printf("Total items (counted): %d\n\n", count)

	// --- List tags ---
	fmt.Println("=== Tags (first 10) ===")
	tags, _, err := client.Tags.List(ctx, lib, zotero.WithLimit(10))
	if err != nil {
		log.Fatalf("Tags.List: %v", err)
	}
	for _, t := range tags {
		num := 0
		if t.Meta != nil {
			num = t.Meta.NumItems
		}
		fmt.Printf("  %s (%d items)\n", t.Tag, num)
	}
	fmt.Println()

	// --- List groups ---
	fmt.Println("=== Groups ===")
	groups, _, err := client.Groups.List(ctx, userID)
	if err != nil {
		log.Fatalf("Groups.List: %v", err)
	}
	for _, g := range groups {
		fmt.Printf("  [%d] %s (%s)\n", g.ID, g.Data.Name, g.Data.Type)
	}
	fmt.Println()

	// --- Schema: available item types ---
	fmt.Println("=== Item Types (first 5) ===")
	itemTypes, _, err := client.Schema.ItemTypes(ctx)
	if err != nil {
		log.Fatalf("Schema.ItemTypes: %v", err)
	}
	for i, it := range itemTypes {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(itemTypes)-5)
			break
		}
		fmt.Printf("  %s (%s)\n", it.ItemType, it.Localized)
	}

	fmt.Println("\nDone.")
}
