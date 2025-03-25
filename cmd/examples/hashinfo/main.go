package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/pixelsquared/go-hashcat"
)

// Example application for listing supported hash types
func main() {
	// Create a new hashcat client with default options
	client, err := hashcat.NewClient()
	if err != nil {
		log.Fatalf("Failed to create hashcat client: %v", err)
	}

	// Get all supported hash types
	fmt.Println("Retrieving supported hash types...")
	hashTypes, err := client.GetSupportedHashes(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving hash types: %v", err)
	}

	// Create a tabwriter for nice formatting
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tCategory")
	fmt.Fprintln(w, "----\t----\t--------")

	// Sort hash types by ID for consistent output
	sort.Slice(hashTypes.HashTypes, func(i, j int) bool {
		return hashTypes.HashTypes[i].ID < hashTypes.HashTypes[j].ID
	})

	// Print information about each hash type
	for _, hashType := range hashTypes.HashTypes {
		fmt.Fprintf(w, "%d\t%s\t%s\n", hashType.ID, hashType.Name, hashType.Category)
	}
	w.Flush()

	// Print summary
	fmt.Printf("\nTotal supported hash types: %d\n", len(hashTypes.HashTypes))

	// Group hash types by category
	categories := make(map[string]int)
	for _, hashType := range hashTypes.HashTypes {
		categories[hashType.Category]++
	}

	// Print category summary
	fmt.Println("\nHash types by category:")
	for category, count := range categories {
		fmt.Printf("  %s: %d\n", category, count)
	}

	// Search example
	fmt.Println("\nSearch example - Finding MD5 hash type:")
	md5Type := hashcat.FindHashTypeByName(hashTypes, "MD5")
	if md5Type != nil {
		fmt.Printf("  MD5 hash type ID: %d\n", md5Type.ID)
	} else {
		fmt.Println("  MD5 hash type not found")
	}
}
