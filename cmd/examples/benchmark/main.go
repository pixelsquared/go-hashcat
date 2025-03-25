package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"text/tabwriter"
	"time"

	"github.com/pixelsquared/go-hashcat"
)

func main() {
	// Parse command line arguments
	hashTypeArg := flag.Int("hash", 0, "Hash type ID to benchmark (default: 0 = MD5)")
	listFlag := flag.Bool("list", false, "List common hash types and exit")
	flag.Parse()

	// Create hashcat client
	client, err := hashcat.NewClient()
	if err != nil {
		log.Fatalf("Failed to create hashcat client: %v", err)
	}

	// If list flag is set, list some common hash types and exit
	if *listFlag {
		printCommonHashTypes(client)
		return
	}

	hashType := *hashTypeArg

	// Get hash type name for display
	hashTypeName := getHashTypeName(client, hashType)

	fmt.Printf("Starting benchmark for hash type %d (%s)...\n", hashType, hashTypeName)
	fmt.Println("This may take a few seconds...")

	// Start benchmark
	startTime := time.Now()
	benchmark, err := client.Benchmark(context.Background(), hashType)
	duration := time.Since(startTime)

	if err != nil {
		log.Fatalf("Benchmark failed: %v", err)
	}

	if len(benchmark.Benchmarks) == 0 {
		log.Fatalf("No benchmark results returned")
	}

	// Print results in a nice table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "\nBenchmark completed in %.2f seconds\n", duration.Seconds())
	fmt.Fprintln(w, "\nHash Mode\tHash Name")
	fmt.Fprintln(w, "---------\t---------")

	for _, b := range benchmark.Benchmarks {
		fmt.Fprintf(w, "%d\t%s\n", b.HashMode, b.HashName)
	}

	fmt.Fprintln(w, "\nDevice ID\tDevice\tSpeed\tUnit\tAccel\tLoops\tThreads\tVector")
	fmt.Fprintln(w, "---------\t------\t-----\t----\t-----\t-----\t-------\t------")

	// Print results for each device
	for _, b := range benchmark.Benchmarks {
		for _, result := range b.DeviceResults {
			fmt.Fprintf(w, "%d\tDevice #%d\t%.2f\t%s\t%d\t%d\t%d\t%d\n",
				result.DeviceID,
				result.DeviceID,
				result.Speed,
				result.SpeedUnit,
				result.Acceleration,
				result.Loops,
				result.Threads,
				result.VectorSize,
			)
		}
	}
	w.Flush()

	// Print summary and advice
	fmt.Println("\nSummary:")

	// Calculate total speed across all devices
	var totalSpeed float64
	var speedUnit string

	if len(benchmark.Benchmarks) > 0 && len(benchmark.Benchmarks[0].DeviceResults) > 0 {
		speedUnit = benchmark.Benchmarks[0].DeviceResults[0].SpeedUnit

		for _, result := range benchmark.Benchmarks[0].DeviceResults {
			totalSpeed += result.Speed
		}
	}

	fmt.Printf("  Total speed: %.2f %s\n", totalSpeed, speedUnit)

	// Provide some context on what the speed means
	estimatedPasswordsPerSecond := totalSpeed

	// Convert to base H/s (from MH/s, GH/s, etc.)
	switch speedUnit {
	case "kH/s":
		estimatedPasswordsPerSecond *= 1000
	case "MH/s":
		estimatedPasswordsPerSecond *= 1000000
	case "GH/s":
		estimatedPasswordsPerSecond *= 1000000000
	}

	fmt.Printf("  Passwords per second: %.0f\n", estimatedPasswordsPerSecond)

	// Provide some time estimates for cracking
	if estimatedPasswordsPerSecond > 0 {
		fmt.Println("\nEstimated time to crack:")

		// 4-char lowercase (26^4)
		lowercase4 := (26.0 * 26.0 * 26.0 * 26.0) / estimatedPasswordsPerSecond
		fmt.Printf("  4-char lowercase password (a-z): %.2f seconds\n", lowercase4)

		// 6-char lowercase (26^6)
		lowercase6 := (26.0 * 26.0 * 26.0 * 26.0 * 26.0 * 26.0) / estimatedPasswordsPerSecond
		formatTimeEstimate(w, "  6-char lowercase password (a-z)", lowercase6)

		// 8-char lowercase (26^8)
		lowercase8 := (26.0 * 26.0 * 26.0 * 26.0 * 26.0 * 26.0 * 26.0 * 26.0) / estimatedPasswordsPerSecond
		formatTimeEstimate(w, "  8-char lowercase password (a-z)", lowercase8)

		// 8-char mixed case + numbers (62^8)
		mixed8 := math.Pow(62.0, 8.0) / estimatedPasswordsPerSecond
		formatTimeEstimate(w, "  8-char alphanumeric password (a-zA-Z0-9)", mixed8)
	}
}

// getHashTypeName returns the name of a hash type given its ID
func getHashTypeName(client hashcat.Client, hashType int) string {
	hashTypes, err := client.GetSupportedHashes(context.Background())
	if err != nil {
		return fmt.Sprintf("Unknown (%v)", err)
	}

	for _, ht := range hashTypes.HashTypes {
		if ht.ID == hashType {
			return ht.Name
		}
	}

	return "Unknown"
}

// printCommonHashTypes prints a list of common hash types
func printCommonHashTypes(client hashcat.Client) {
	// Get all hash types
	hashTypes, err := client.GetSupportedHashes(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving hash types: %v", err)
	}

	// Define some common hash types to list
	commonHashIDs := []int{0, 100, 1000, 1800, 3000, 5500, 5600, 1400, 1700}

	fmt.Println("Common hash types for benchmarking:")
	fmt.Println("----------------------------------")

	for _, id := range commonHashIDs {
		for _, ht := range hashTypes.HashTypes {
			if ht.ID == id {
				fmt.Printf("%d = %s\n", id, ht.Name)
				break
			}
		}
	}

	fmt.Println("\nTo benchmark a specific hash type:")
	fmt.Println("  go run main.go -hash <ID>")
	fmt.Println("\nExample:")
	fmt.Println("  go run main.go -hash 1000  # Benchmark NTLM hashes")
}

// formatTimeEstimate formats a time estimate in human-readable form
func formatTimeEstimate(w *tabwriter.Writer, label string, seconds float64) {
	if seconds < 60 {
		fmt.Printf("%s: %.2f seconds\n", label, seconds)
		return
	}

	minutes := seconds / 60
	if minutes < 60 {
		fmt.Printf("%s: %.2f minutes\n", label, minutes)
		return
	}

	hours := minutes / 60
	if hours < 24 {
		fmt.Printf("%s: %.2f hours\n", label, hours)
		return
	}

	days := hours / 24
	if days < 365 {
		fmt.Printf("%s: %.2f days\n", label, days)
		return
	}

	years := days / 365
	fmt.Printf("%s: %.2f years\n", label, years)
}
