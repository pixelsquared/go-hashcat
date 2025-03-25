package models

// Benchmark represents the performance benchmark of hashcat for a specific hash type
type Benchmark struct {
	HashMode      int               `json:"hash_mode"`
	HashName      string            `json:"hash_name"`
	DeviceResults []BenchmarkResult `json:"device_results"`
}

// BenchmarkResult represents the benchmark result for a specific device
type BenchmarkResult struct {
	DeviceID     int     `json:"device_id"`
	Speed        float64 `json:"speed"`
	SpeedUnit    string  `json:"speed_unit"`
	TimePerHash  float64 `json:"time_per_hash_ms"`
	Acceleration int     `json:"acceleration"`
	Loops        int     `json:"loops"`
	Threads      int     `json:"threads"`
	VectorSize   int     `json:"vector_size"`
}

// BenchmarkSummary represents the summarized benchmark results across all devices
type BenchmarkSummary struct {
	TotalSpeed     float64 `json:"total_speed"`
	SpeedUnit      string  `json:"speed_unit"`
	AvgTimePerHash float64 `json:"avg_time_per_hash_ms"`
}

// HashcatBenchmarkResponse represents the full response from hashcat's benchmark command
type HashcatBenchmarkResponse struct {
	Benchmarks []Benchmark      `json:"benchmarks"`
	Summary    BenchmarkSummary `json:"summary,omitempty"`
}
