package hashcat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pixelsquared/go-hashcat/models"
)

// ParseDeviceOutput parses the output from hashcat --backend-info command
func ParseDeviceOutput(output string) (*models.DeviceList, error) {
	lines := strings.Split(output, "\n")
	var platforms []models.Platform
	var currentPlatform *models.Platform
	var currentDevice *models.Device

	// Regular expressions for parsing
	rePlatform := regexp.MustCompile(`(?:^|\n)OpenCL Platform ID #(\d+)`)
	reDevice := regexp.MustCompile(`(?:^|\n)\s*Backend Device ID #(\d+)`)

	// Property patterns
	platformVendorRe := regexp.MustCompile(`^\s+Vendor\.\.: (.+)$`)
	platformNameRe := regexp.MustCompile(`^\s+Name\....: (.+)$`)
	platformVersionRe := regexp.MustCompile(`^\s+Version\.: (.+)$`)

	// Device properties
	deviceTypeRe := regexp.MustCompile(`^\s+Type\.\.\.\.\.\.\.\.\.\.\.: (.+)$`)
	deviceVendorIDRe := regexp.MustCompile(`^\s+Vendor\.ID\.\.\.\.\.\.: (\d+)$`)
	deviceVendorRe := regexp.MustCompile(`^\s+Vendor\.\.\.\.\.\.\.\.\.: (.+)$`)
	deviceNameRe := regexp.MustCompile(`^\s+Name\.\.\.\.\.\.\.\.\.\.\.: (.+)$`)
	deviceVersionRe := regexp.MustCompile(`^\s+Version\.\.\.\.\.\.\.\.: (.+)$`)
	deviceProcessorsRe := regexp.MustCompile(`^\s+Processor\(s\)\.\.\.: (\d+)$`)
	deviceClockRe := regexp.MustCompile(`^\s+Clock\.\.\.\.\.\.\.\.\.: (\d+)$`)
	deviceMemTotalRe := regexp.MustCompile(`^\s+Memory\.Total\.\.\.: (\d+)`)
	deviceMemFreeRe := regexp.MustCompile(`^\s+Memory\.Free\.\.\.\.: (\d+)`)
	deviceLocalMemRe := regexp.MustCompile(`^\s+Local\.Memory\.\.\.: (\d+)`)
	deviceOpenCLVersionRe := regexp.MustCompile(`^\s+OpenCL\.Version\.: (.+)$`)
	deviceDriverVersionRe := regexp.MustCompile(`^\s+Driver\.Version\.: (.+)$`)

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check for platform
		platformMatch := rePlatform.FindStringSubmatch(line)
		if len(platformMatch) > 1 {
			// Found a platform, create new platform
			platformID, _ := strconv.Atoi(platformMatch[1])
			currentPlatform = &models.Platform{
				ID:      platformID,
				Devices: []models.Device{},
			}
			currentDevice = nil // Reset current device when changing platform
			platforms = append(platforms, *currentPlatform)
			continue
		}

		// Check for device
		deviceMatch := reDevice.FindStringSubmatch(line)
		if len(deviceMatch) > 1 && currentPlatform != nil {
			// Found a device, create new device in the latest platform
			deviceID, _ := strconv.Atoi(deviceMatch[1])
			device := models.Device{
				ID: deviceID,
			}
			platformIdx := len(platforms) - 1
			platforms[platformIdx].Devices = append(platforms[platformIdx].Devices, device)
			// Get a reference to the newly added device
			currentDevice = &platforms[platformIdx].Devices[len(platforms[platformIdx].Devices)-1]
			continue
		}

		// Check for platform properties
		if currentPlatform != nil && currentDevice == nil {
			if match := platformVendorRe.FindStringSubmatch(line); len(match) > 1 {
				platforms[len(platforms)-1].Vendor = match[1]
			} else if match := platformNameRe.FindStringSubmatch(line); len(match) > 1 {
				platforms[len(platforms)-1].Name = match[1]
			} else if match := platformVersionRe.FindStringSubmatch(line); len(match) > 1 {
				platforms[len(platforms)-1].Version = match[1]
			}
		}

		// Check for device properties
		if currentDevice != nil && currentPlatform != nil {
			platformIdx := len(platforms) - 1
			deviceIdx := len(platforms[platformIdx].Devices) - 1

			// Make sure we have valid indices before proceeding
			if platformIdx >= 0 && deviceIdx >= 0 {
				if match := deviceTypeRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].Type = match[1]
				} else if match := deviceVendorIDRe.FindStringSubmatch(line); len(match) > 1 {
					vendorID, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].VendorID = vendorID
				} else if match := deviceVendorRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].Vendor = match[1]
				} else if match := deviceNameRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].Name = match[1]
				} else if match := deviceVersionRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].Version = match[1]
				} else if match := deviceProcessorsRe.FindStringSubmatch(line); len(match) > 1 {
					processors, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].Processors = processors
				} else if match := deviceClockRe.FindStringSubmatch(line); len(match) > 1 {
					clock, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].ClockMHz = clock
				} else if match := deviceMemTotalRe.FindStringSubmatch(line); len(match) > 1 {
					memTotal, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].MemoryTotal = memTotal
				} else if match := deviceMemFreeRe.FindStringSubmatch(line); len(match) > 1 {
					memFree, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].MemoryFree = memFree
				} else if match := deviceLocalMemRe.FindStringSubmatch(line); len(match) > 1 {
					localMem, _ := strconv.Atoi(match[1])
					platforms[platformIdx].Devices[deviceIdx].LocalMemory = localMem
				} else if match := deviceOpenCLVersionRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].OpenCLVersion = match[1]
				} else if match := deviceDriverVersionRe.FindStringSubmatch(line); len(match) > 1 {
					platforms[platformIdx].Devices[deviceIdx].DriverVersion = match[1]
				}
			}
		}
	}

	if len(platforms) == 0 {
		return nil, fmt.Errorf("no OpenCL platforms found in output")
	}

	return &models.DeviceList{
		Platforms: platforms,
	}, nil
}

// FindDeviceByID finds a device with the given ID across all platforms
func FindDeviceByID(devices *models.DeviceList, id int) *models.Device {
	for _, platform := range devices.Platforms {
		for _, device := range platform.Devices {
			if device.ID == id {
				return &device
			}
		}
	}
	return nil
}

// FindDeviceByName finds a device with the given name (case insensitive)
func FindDeviceByName(devices *models.DeviceList, name string) *models.Device {
	lowerName := strings.ToLower(name)
	for _, platform := range devices.Platforms {
		for _, device := range platform.Devices {
			if strings.Contains(strings.ToLower(device.Name), lowerName) {
				return &device
			}
		}
	}
	return nil
}

// FindDevicesByType returns all devices of a specific type (e.g., "CPU" or "GPU")
func FindDevicesByType(devices *models.DeviceList, deviceType string) []models.Device {
	var result []models.Device
	for _, platform := range devices.Platforms {
		for _, device := range platform.Devices {
			if strings.EqualFold(device.Type, deviceType) {
				result = append(result, device)
			}
		}
	}
	return result
}

// FindPlatformByID returns a platform with the given ID
func FindPlatformByID(devices *models.DeviceList, id int) *models.Platform {
	for _, platform := range devices.Platforms {
		if platform.ID == id {
			return &platform
		}
	}
	return nil
}
