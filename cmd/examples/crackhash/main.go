package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pixelsquared/go-hashcat"
	"github.com/pixelsquared/go-hashcat/models"
)

func main() {
	// Parse command line arguments
	hashArg := flag.String("hash", "5f4dcc3b5aa765d61d8327deb882cf99", "Hash to crack (default: MD5 of 'password')")
	hashTypeArg := flag.Int("type", 0, "Hash type ID (default: 0 = MD5)")
	attackModeArg := flag.Int("attack", 3, "Attack mode (0=dict, 3=mask, etc.)")
	maskArg := flag.String("mask", "?a?a?a?a?a?a?a?a", "Mask or dictionary file path")
	verboseArg := flag.Bool("verbose", false, "Show detailed progress information")
	flag.Parse()

	// Print banner
	fmt.Println("====================================")
	fmt.Println("go-hashcat Example: Hash Cracking")
	fmt.Println("====================================")
	fmt.Println()

	// Create hashcat client
	client, err := hashcat.NewClient()
	if err != nil {
		log.Fatalf("Failed to create hashcat client: %v", err)
	}

	// Display parameters
	fmt.Printf("Hash:       %s\n", *hashArg)
	fmt.Printf("Hash Type:  %d\n", *hashTypeArg)
	fmt.Printf("Attack:     %d\n", *attackModeArg)
	fmt.Printf("Mask/Dict:  %s\n", *maskArg)
	fmt.Println()

	// Set up context with cancelation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, stopping...")
		cancel()
	}()

	fmt.Println("Starting cracking session...")
	startTime := time.Now()

	// Create options for the crack session
	options := &hashcat.CrackOptions{
		HashType:        *hashTypeArg,
		AttackMode:      *attackModeArg,
		Mask:            *maskArg,
		OptimizedKernel: true,
		Workload:        2, // Default workload
	}

	// Create a new crack session
	session, err := client.NewCrackSession(ctx, *hashArg, options)
	if err != nil {
		log.Fatalf("Failed to create crack session: %v", err)
	}

	// Get progress channel
	progressChan := session.Progress()

	// Start the cracking process
	if err := session.Start(); err != nil {
		log.Fatalf("Failed to start cracking: %v", err)
	}

	fmt.Println("Cracking in progress. Press Ctrl+C to stop.")
	fmt.Println()

	// Display header for progress
	if *verboseArg {
		fmt.Printf("%-10s | %-10s | %-20s | %-15s | %s\n",
			"Status", "Progress", "Speed", "Elapsed", "ETA")
		fmt.Println(strings.Repeat("-", 80))
	}

	// Monitor progress
	lastProgressUpdate := time.Now()
	lastStatus := models.StatusUnknown

	for progress := range progressChan {
		stats := progress.CalculateStats()

		// Get total speed across all devices
		var totalSpeed int64

		for _, device := range progress.Devices {
			totalSpeed += device.Speed
		}

		// Format speed with appropriate unit
		speedFormatted := formatSpeed(totalSpeed)

		// Status message
		var statusMsg string
		switch progress.Status {
		case models.StatusInit:
			statusMsg = "Initializing"
		case models.StatusRunning:
			statusMsg = "Running"
		case models.StatusExhausted:
			statusMsg = "Exhausted"
		case models.StatusCracked:
			statusMsg = "Cracked!"
		case models.StatusAborted:
			statusMsg = "Aborted"
		case models.StatusQuit:
			statusMsg = "Quit"
		case models.StatusPaused:
			statusMsg = "Paused"
		default:
			statusMsg = "Unknown"
		}

		// Print progress updates less frequently to avoid spam
		now := time.Now()
		if *verboseArg && (now.Sub(lastProgressUpdate) > 1*time.Second || progress.Status != lastStatus) {
			fmt.Printf("%-10s | %6.2f%%    | %-20s | %-15s | %s\n",
				statusMsg,
				stats.PercentComplete,
				speedFormatted,
				formatDuration(stats.ElapsedTime),
				formatDuration(stats.EstimatedRemaining),
			)
			lastProgressUpdate = now
		} else if !*verboseArg {
			// Simple progress bar for non-verbose mode
			fmt.Printf("\r[%-30s] %6.2f%% | %s | %s remaining",
				progressBar(stats.PercentComplete, 30),
				stats.PercentComplete,
				speedFormatted,
				formatDuration(stats.EstimatedRemaining),
			)
		}

		// Update last status
		lastStatus = progress.Status

		// If cracking completed or was stopped, break the loop
		if progress.Status == models.StatusCracked ||
			progress.Status == models.StatusExhausted ||
			progress.Status == models.StatusAborted {
			break
		}
	}

	// Add newline after progress output
	fmt.Println()

	// Wait for completion and get results
	if err := session.Wait(); err != nil {
		log.Fatalf("Error during cracking: %v", err)
	}

	results, err := session.Results()
	if err != nil {
		log.Fatalf("Error getting results: %v", err)
	}

	// Print final results
	duration := time.Since(startTime)
	fmt.Println("\nCracking session completed in", formatDuration(duration))

	if len(results) > 0 {
		fmt.Println("\nCracked hashes:")
		fmt.Println("==============")

		for i, result := range results {
			fmt.Printf("%d. %s = %s\n", i+1, result.Hash, result.Password)
		}

		fmt.Printf("\nSuccessfully cracked %d hash(es)\n", len(results))
	} else {
		fmt.Println("\nNo hashes were cracked.")
	}
}

// formatSpeed formats a speed value with appropriate units
func formatSpeed(speed int64) string {
	if speed < 1000 {
		return fmt.Sprintf("%d H/s", speed)
	} else if speed < 1000000 {
		return fmt.Sprintf("%.2f kH/s", float64(speed)/1000)
	} else if speed < 1000000000 {
		return fmt.Sprintf("%.2f MH/s", float64(speed)/1000000)
	} else {
		return fmt.Sprintf("%.2f GH/s", float64(speed)/1000000000)
	}
}

// formatDuration formats a duration in a human-readable form
func formatDuration(d time.Duration) string {
	// Round to seconds
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	} else {
		return fmt.Sprintf("%ds", s)
	}
}

// progressBar generates a simple ASCII progress bar
func progressBar(percent float64, width int) string {
	completed := int(percent / 100 * float64(width))
	if completed > width {
		completed = width
	}

	bar := strings.Repeat("=", completed)
	if completed < width {
		bar += ">"
		bar += strings.Repeat(" ", width-completed-1)
	}

	return bar
}
