package hashcat

import (
	"context"
	"fmt"

	"github.com/pixelsquared/go-hashcat/models"
)

// BenchmarkAll performs benchmarks for all supported hash types
func (c *HashcatClient) BenchmarkAll(ctx context.Context) (*models.HashcatBenchmarkResponse, error) {
	// Get all supported hash types first
	hashes, err := c.GetSupportedHashes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported hash types: %w", err)
	}

	response := &models.HashcatBenchmarkResponse{
		Benchmarks: []models.Benchmark{},
	}

	// Benchmark each hash type
	for _, hashType := range hashes.HashTypes {
		benchmarkResult, err := c.Benchmark(ctx, hashType.ID)
		if err != nil {
			// Log error but continue with other hash types
			continue
		}

		if len(benchmarkResult.Benchmarks) > 0 {
			response.Benchmarks = append(response.Benchmarks, benchmarkResult.Benchmarks...)
		}
	}

	// Calculate overall summary
	var totalSpeed float64
	var totalDevices int
	var totalTimePerHash float64

	// Use a common speed unit for the summary (MH/s)
	speedUnit := "MH/s"

	for _, benchmark := range response.Benchmarks {
		for _, result := range benchmark.DeviceResults {
			// Normalize speed to MH/s for consistent reporting
			normalizedSpeed := result.Speed
			switch result.SpeedUnit {
			case "GH/s":
				normalizedSpeed *= 1000
			case "kH/s":
				normalizedSpeed /= 1000
			case "H/s":
				normalizedSpeed /= 1000000
			}

			totalSpeed += normalizedSpeed
			totalTimePerHash += result.TimePerHash
			totalDevices++
		}
	}

	// Calculate average time per hash
	avgTimePerHash := 0.0
	if totalDevices > 0 {
		avgTimePerHash = totalTimePerHash / float64(totalDevices)
	}

	response.Summary = models.BenchmarkSummary{
		TotalSpeed:     totalSpeed,
		SpeedUnit:      speedUnit,
		AvgTimePerHash: avgTimePerHash,
	}

	return response, nil
}
