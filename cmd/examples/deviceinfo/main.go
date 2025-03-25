package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pixelsquared/go-hashcat"
)

// Example application for listing available devices
func main() {
	// Create a new hashcat client with default options
	client, err := hashcat.NewClient()
	if err != nil {
		log.Fatalf("Failed to create hashcat client: %v", err)
	}

	// Get information about available devices
	fmt.Println("Retrieving device information...")
	devices, err := client.GetDevices(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving device information: %v", err)
	}

	if len(devices.Platforms) == 0 {
		fmt.Println("No compatible platforms found.")
		return
	}

	// Print platform and device information
	for i, platform := range devices.Platforms {
		fmt.Printf("\nPlatform #%d: %s\n", platform.ID, platform.Name)
		fmt.Printf("  Vendor:  %s\n", platform.Vendor)
		fmt.Printf("  Version: %s\n", platform.Version)

		if len(platform.Devices) == 0 {
			fmt.Println("  No devices found for this platform.")
			continue
		}

		fmt.Println("\n  Devices:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "    ID\tName\tType\tProcessors\tClock\tMemory")
		fmt.Fprintln(w, "    --\t----\t----\t----------\t-----\t------")

		for _, device := range platform.Devices {
			fmt.Fprintf(w, "    %d\t%s\t%s\t%d\t%d MHz\t%d MB\n",
				device.ID,
				device.Name,
				device.Type,
				device.Processors,
				device.ClockMHz,
				device.MemoryTotal)
		}
		w.Flush()

		// Print device details for the first device in each platform
		if len(platform.Devices) > 0 {
			device := platform.Devices[0]
			fmt.Printf("\n  Detailed information for %s:\n", device.Name)
			fmt.Printf("    Vendor ID:       %d\n", device.VendorID)
			fmt.Printf("    Vendor:          %s\n", device.Vendor)
			fmt.Printf("    OpenCL Version:  %s\n", device.OpenCLVersion)
			fmt.Printf("    Driver Version:  %s\n", device.DriverVersion)
			fmt.Printf("    Local Memory:    %d KB\n", device.LocalMemory)
			fmt.Printf("    Memory Free:     %d MB\n", device.MemoryFree)
		}

		// Only add separator if there are more platforms
		if i < len(devices.Platforms)-1 {
			fmt.Println("\n" + strings.Repeat("-", 80))
		}
	}

	// Print summary statistics
	fmt.Println("\nSummary:")
	totalDevices := 0
	gpuDevices := 0
	cpuDevices := 0

	for _, platform := range devices.Platforms {
		totalDevices += len(platform.Devices)
		for _, device := range platform.Devices {
			if strings.Contains(strings.ToLower(device.Type), "gpu") {
				gpuDevices++
			} else if strings.Contains(strings.ToLower(device.Type), "cpu") {
				cpuDevices++
			}
		}
	}

	fmt.Printf("  Total platforms: %d\n", len(devices.Platforms))
	fmt.Printf("  Total devices:   %d\n", totalDevices)
	fmt.Printf("  GPU devices:     %d\n", gpuDevices)
	fmt.Printf("  CPU devices:     %d\n", cpuDevices)

	// Example of using helper functions
	fmt.Println("\nHelper function examples:")

	// Find devices by type
	gpus := hashcat.FindDevicesByType(devices, "GPU")
	fmt.Printf("  Found %d GPU devices\n", len(gpus))

	// Find a specific device by name (example: search for "NVIDIA")
	nvidiaDevice := hashcat.FindDeviceByName(devices, "NVIDIA")
	if nvidiaDevice != nil {
		fmt.Printf("  Found NVIDIA device: %s\n", nvidiaDevice.Name)
	} else {
		fmt.Println("  No NVIDIA device found")
	}
}
