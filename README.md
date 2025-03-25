# go-hashcat

[![Go Reference](https://pkg.go.dev/badge/github.com/pixelsquared/go-hashcat.svg)](https://pkg.go.dev/github.com/pixelsquared/go-hashcat)
[![Go Report Card](https://goreportcard.com/badge/github.com/pixelsquared/go-hashcat)](https://goreportcard.com/report/github.com/pixelsquared/go-hashcat)

A Go library for interacting with hashcat.

## Overview

The go-hashcat library provides a Go interface to the hashcat password recovery tool, making it easy to integrate hashcat's powerful functionality into Go applications. The library wraps the hashcat command-line tool, providing a simple and idiomatic Go API that abstracts away the complexity of direct command execution.

Key features:
- Hardware information discovery (CPU, GPU capabilities)
- Hash type identification and information
- Benchmarking password cracking performance
- Password cracking with various attack modes
- Real-time progress monitoring through channels
- Structured data types for all hashcat operations

## Requirements

- Go 1.16 or higher
- Hashcat 6.0.0 or higher installed and available in PATH (or configured via options)
- Compatible hardware for GPU acceleration (optional, but recommended for performance)

### Supported Platforms

- Linux
- Windows
- macOS (limited GPU support)

## Installation

```bash
go get github.com/pixelsquared/go-hashcat
```

### Installing Hashcat

You'll need to install the hashcat binary separately:

#### Linux
```bash
sudo apt-get install hashcat
```

#### Windows
Download and install from the [official hashcat website](https://hashcat.net/hashcat/).

#### macOS
```bash
brew install hashcat
```

## Basic Usage

### Creating a Client

```go
import (
    "context"
    "github.com/pixelsquared/go-hashcat"
)

// Create a new client with default options
client, err := hashcat.NewClient()
if err != nil {
    log.Fatalf("Failed to create hashcat client: %v", err)
}

// Or with custom options
client, err := hashcat.NewClient(
    hashcat.WithBinaryPath("/path/to/hashcat"),
    hashcat.WithOutputDir("/path/to/output"),
)
if err != nil {
    log.Fatalf("Failed to create hashcat client: %v", err)
}
```

### Getting Hardware Information

```go
// Get information about available devices
devices, err := client.GetDevices(context.Background())
if err != nil {
    log.Fatalf("Error getting devices: %v", err)
}

// Print device information
for _, platform := range devices.Platforms {
    fmt.Printf("Platform: %s\n", platform.Name)
    
    for _, device := range platform.Devices {
        fmt.Printf("  Device: %s (%s)\n", device.Name, device.Type)
        fmt.Printf("  Memory: %d MB\n", device.MemoryTotal)
        fmt.Printf("  Clock:  %d MHz\n", device.ClockMHz)
    }
}
```

### Listing Supported Hash Types

```go
// Get all supported hash types
hashTypes, err := client.GetSupportedHashes(context.Background())
if err != nil {
    log.Fatalf("Error getting hash types: %v", err)
}

// Print hash type information
for _, hashType := range hashTypes.HashTypes {
    fmt.Printf("ID: %d, Name: %s, Category: %s\n", 
        hashType.ID, hashType.Name, hashType.Category)
}

// Find a specific hash type by name
md5Type := hashcat.FindHashTypeByName(hashTypes, "MD5")
if md5Type != nil {
    fmt.Printf("MD5 hash type ID: %d\n", md5Type.ID)
}
```

### Running Benchmarks

```go
// Run benchmark for MD5 (hash type 0)
benchmark, err := client.Benchmark(context.Background(), 0)
if err != nil {
    log.Fatalf("Benchmark failed: %v", err)
}

// Print benchmark results
for _, b := range benchmark.Benchmarks {
    fmt.Printf("Hash: %s (Mode: %d)\n", b.HashName, b.HashMode)
    
    for _, result := range b.DeviceResults {
        fmt.Printf("  Device #%d: %.2f %s\n", 
            result.DeviceID, result.Speed, result.SpeedUnit)
    }
}
```

### Cracking a Hash

```go
// Basic hash cracking
options := &hashcat.CrackOptions{
    HashType:   0,       // MD5
    AttackMode: 3,       // Mask attack
    Mask:       "?a?a?a?a?a?a?a?a", // 8 chars, all charsets
}

session, err := client.NewCrackSession(context.Background(), "5f4dcc3b5aa765d61d8327deb882cf99", options)
if err != nil {
    log.Fatalf("Failed to create crack session: %v", err)
}

// Start cracking
if err := session.Start(); err != nil {
    log.Fatalf("Failed to start cracking: %v", err)
}

// Monitor progress
for progress := range session.Progress() {
    fmt.Printf("\rProgress: %.2f%%, Speed: %d H/s",
        progress.CalculateStats().PercentComplete,
        progress.Devices[0].Speed)
        
    if progress.Status == models.StatusCracked {
        break
    }
}

// Get results
results, err := session.Results()
if err != nil {
    log.Fatalf("Error getting results: %v", err)
}

for _, result := range results {
    fmt.Printf("\nCracked: %s = %s\n", result.Hash, result.Password)
}
```

## Advanced Usage

### Dictionary Attack

```go
options := &hashcat.CrackOptions{
    HashType:   0,       // MD5
    AttackMode: 0,       // Dictionary attack
    Mask:       "/path/to/wordlist.txt",
    RulesFile:  "/path/to/rules.rule", // Optional
}

session, err := client.NewCrackSession(context.Background(), "5f4dcc3b5aa765d61d8327deb882cf99", options)
```

### Brute Force with Custom Character Sets

```go
options := &hashcat.CrackOptions{
    HashType:    0,                 // MD5
    AttackMode:  3,                 // Mask attack
    Mask:        "?l?l?l?l?d?d?d?d", // 4 lowercase + 4 digits
    CustomCharset1: "abcdefghijklmnopqrstuvwxyz", // Custom set for ?1
}
```

### Session Management

```go
// Create a session with a specific name for resuming later
options := &hashcat.CrackOptions{
    HashType:     0,                // MD5
    AttackMode:   3,                // Mask attack
    Mask:         "?a?a?a?a?a?a?a?a",
    SessionName:  "my-crack-session",
}

// Pause a running session
if err := session.Pause(); err != nil {
    log.Printf("Error pausing session: %v", err)
}

// Resume a session
if err := session.Resume(); err != nil {
    log.Printf("Error resuming session: %v", err)
}

// Shutdown gracefully
if err := session.Shutdown(context.Background()); err != nil {
    log.Printf("Error shutting down: %v", err)
}
```

## API Documentation

### Client Interface

The main interface for interacting with hashcat:

```go
type Client interface {
    // Get information about available devices
    GetDevices(ctx context.Context) (*models.Devices, error)
    
    // Get a list of supported hash types
    GetSupportedHashes(ctx context.Context) (*models.HashTypes, error)
    
    // Run benchmark for a specific hash type
    Benchmark(ctx context.Context, hashType int) (*models.BenchmarkResults, error)
    
    // Create a new cracking session
    NewCrackSession(ctx context.Context, hash string, options *CrackOptions) (CrackSession, error)
    
    // Crack a hash directly (simplified interface)
    Crack(ctx context.Context, hash string, hashType, attackMode int, mask string) (<-chan models.Progress, error)
}
```

### CrackSession Interface

Interface for managing an individual cracking session:

```go
type CrackSession interface {
    // Start the cracking process
    Start() error
    
    // Get a channel for monitoring progress
    Progress() <-chan models.Progress
    
    // Wait for the session to complete
    Wait() error
    
    // Get results after completion
    Results() ([]models.CrackedHash, error)
    
    // Pause a running session
    Pause() error
    
    // Resume a paused session
    Resume() error
    
    // Stop the session
    Stop() error
    
    // Shutdown gracefully with timeout
    Shutdown(ctx context.Context) error
}
```

### Key Data Models

#### Device Information

```go
type Devices struct {
    Platforms []Platform
}

type Platform struct {
    ID       int
    Name     string
    Vendor   string
    Version  string
    Devices  []Device
}

type Device struct {
    ID            int
    Name          string
    Type          string
    Vendor        string
    VendorID      int
    Processors    int
    ClockMHz      int
    MemoryTotal   int
    MemoryFree    int
    LocalMemory   int
    OpenCLVersion string
    DriverVersion string
}
```

#### Hash Types

```go
type HashTypes struct {
    HashTypes []HashType
}

type HashType struct {
    ID       int
    Name     string
    Category string
}
```

#### Benchmark Results

```go
type BenchmarkResults struct {
    Benchmarks []Benchmark
}

type Benchmark struct {
    HashMode      int
    HashName      string
    DeviceResults []DeviceResult
}

type DeviceResult struct {
    DeviceID      int
    Speed         float64
    SpeedUnit     string
    Acceleration  int
    Loops         int
    Threads       int
    VectorSize    int
}
```

#### Progress Information

```go
type Progress struct {
    Status     Status
    Progress   int
    Speed      int64
    Recovered  int
    Remaining  int
    Rejected   int
    Restored   int
    Salts      Salt
    Devices    []DeviceProgress
    TimeStart  time.Time
    TimeStop   time.Time
    Checkpoint time.Duration
}

type Salt struct {
    Current  int
    Total    int
    Recovered int
}

type DeviceProgress struct {
    ID        int
    Speed     int64
    Recovered int
    Rejected  int
}

// Stats provides analytics derived from progress data
type Stats struct {
    PercentComplete    float64
    ElapsedTime        time.Duration
    EstimatedRemaining time.Duration
    EstimatedTotal     time.Duration
}
```

#### Cracked Hashes

```go
type CrackedHash struct {
    Hash     string
    Password string
    Type     int
}
```

## Configuration Options

The library supports various configuration options through functional options pattern:

```go
// Specify path to hashcat binary
hashcat.WithBinaryPath("/path/to/hashcat")

// Specify directory for hashcat output files
hashcat.WithOutputDir("/path/to/output")

// Set limits on CPU usage
hashcat.WithCPUAffinity("1,2,3")

// Configure GPU devices to use
hashcat.WithDevices("1,2") // Use devices 1 and 2

// Configure workload profile (0-4)
hashcat.WithWorkload(2)

// Set advanced optimization options
hashcat.WithOptimizedKernel(true)
hashcat.WithSelfTestDisabled(true)
```

## Error Handling

The library provides specialized error types for common scenarios:

```go
// Check for specific error types
switch err := err.(type) {
case *hashcat.HashcatNotFoundError:
    log.Fatalf("Hashcat binary not found: %v", err)
case *hashcat.HashcatExecutionError:
    log.Fatalf("Hashcat execution failed: %v", err.Error())
    log.Printf("Command output: %s", err.Output)
case *hashcat.SessionError:
    log.Fatalf("Session error: %v", err)
default:
    log.Fatalf("Unknown error: %v", err)
}
```

## Examples

See the `cmd/examples` directory for complete working examples:

- `hashinfo`: List supported hash types
- `deviceinfo`: Display information about available devices
- `benchmark`: Run performance benchmarks
- `crackhash`: Crack a hash with real-time progress monitoring
