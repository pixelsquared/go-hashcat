// Package hashcat provides a Go interface to hashcat, the world's fastest password recovery tool.
//
// This package wraps the hashcat command-line tool, providing a simple and idiomatic Go API
// for accessing hashcat's functionality. It handles executing the hashcat binary, parsing its
// output, and providing structured data types for working with hashcat's features.
//
// # Overview
//
// The hashcat package is designed to make it easy to integrate hashcat's functionality into Go
// applications. It provides a client interface for interacting with hashcat and abstracts away
// the details of executing commands and parsing output. The package includes support for:
//
// - Getting information about available devices (CPU, GPU)
// - Listing supported hash types
// - Running benchmarks
// - Cracking passwords with various attack modes
// - Monitoring progress in real-time
//
// # Basic Usage
//
// Create a new hashcat client:
//
//	client, err := hashcat.NewClient()
//	if err != nil {
//	    log.Fatalf("Failed to create hashcat client: %v", err)
//	}
//
// Get information about available devices:
//
//	devices, err := client.GetDevices(context.Background())
//	if err != nil {
//	    log.Fatalf("Error getting devices: %v", err)
//	}
//
// List supported hash types:
//
//	hashTypes, err := client.GetSupportedHashes(context.Background())
//	if err != nil {
//	    log.Fatalf("Error getting hash types: %v", err)
//	}
//
// Run a benchmark:
//
//	benchmark, err := client.Benchmark(context.Background(), 0) // MD5
//	if err != nil {
//	    log.Fatalf("Benchmark failed: %v", err)
//	}
//
// Crack a hash:
//
//	progressChan, err := client.Crack(
//	    context.Background(),
//	    "5f4dcc3b5aa765d61d8327deb882cf99", // "password" in MD5
//	    0,                                   // Hash type: MD5
//	    3,                                   // Attack mode: Mask
//	    "?a?a?a?a?a?a?a?a",                 // Mask: 8 chars, all character sets
//	)
//
// # Design Principles
//
// The hashcat package follows these design principles:
//
// 1. Provide a simple, idiomatic Go API that abstracts away the complexity of hashcat.
//
// 2. Use Go's concurrency features to handle asynchronous operations like progress monitoring.
//
// 3. Follow a modular design with clear interfaces and separation of concerns.
//
// 4. Provide structured data types for all hashcat operations and results.
//
// 5. Use context.Context for cancelation and timeouts to allow proper resource management.
//
// 6. Support configuration through functional options for flexibility.
//
// # Implementation Details
//
// The hashcat package uses the following approach:
//
// - The Client interface defines the main API for interacting with hashcat.
//
// - The HashcatClient type implements the Client interface.
//
// - Configuration is handled via functional options (WithBinaryPath, WithOutputDir, etc.).
//
// - Hashcat commands are executed as separate processes, with output parsed into structured data.
//
// - Progress monitoring is implemented using channels to provide real-time updates.
//
// - Sessions are used to manage cracking operations, including starting, stopping, and monitoring.
//
// - All operations support context.Context for cancelation and timeouts.
//
// # Concurrency
//
// The hashcat package is designed to be safe for concurrent use. Each cracking session
// operates independently, with proper synchronization using mutexes and channels.
//
// Progress updates are sent through channels, allowing for non-blocking monitoring of
// cracking operations.
//
// # Error Handling
//
// Errors are wrapped with context to provide more information about the operation that failed.
// Custom error types like HashcatError provide details about the command that failed and
// the output from hashcat.
//
// # Examples
//
// See the examples directory for working examples of how to use the hashcat package.
package hashcat
