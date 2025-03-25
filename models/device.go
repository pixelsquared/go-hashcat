package models

// Platform represents an OpenCL platform in hashcat
type Platform struct {
	ID      int      `json:"id"`
	Vendor  string   `json:"vendor"`
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Devices []Device `json:"devices,omitempty"`
}

// Device represents a backend device in hashcat
type Device struct {
	ID            int    `json:"id"`
	Type          string `json:"type"`
	VendorID      int    `json:"vendor_id"`
	Vendor        string `json:"vendor"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Processors    int    `json:"processors"`
	ClockMHz      int    `json:"clock_mhz"`
	MemoryTotal   int    `json:"memory_total_mb"`
	MemoryFree    int    `json:"memory_free_mb"`
	LocalMemory   int    `json:"local_memory_kb"`
	OpenCLVersion string `json:"opencl_version"`
	DriverVersion string `json:"driver_version"`
}

// DeviceList represents the response from hashcat when listing backend devices
type DeviceList struct {
	Platforms []Platform `json:"platforms"`
}
