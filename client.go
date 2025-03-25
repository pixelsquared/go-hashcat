package hashcat

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pixelsquared/go-hashcat/models"
)

// Client is the interface for interacting with hashcat
type Client interface {
	// GetDevices returns information about available devices
	GetDevices(ctx context.Context) (*models.DeviceList, error)

	// GetSupportedHashes returns information about supported hash types
	GetSupportedHashes(ctx context.Context) (*models.HashcatSupportedHashes, error)

	// Benchmark performs a benchmark for the given hash type
	Benchmark(ctx context.Context, hashType int) (*models.HashcatBenchmarkResponse, error)

	// Crack attempts to crack the provided hash using the specified attack mode and options
	Crack(ctx context.Context, hash string, hashType int, attackMode int, mask string) (<-chan *models.Progress, error)

	// CrackFile attempts to crack hashes in the specified file
	CrackFile(ctx context.Context, hashFile *models.HashFile, attackMode int, mask string) (<-chan *models.Progress, error)

	// Stop interrupts a running cracking session
	Stop(ctx context.Context) error
}

// HashcatClient is the concrete implementation of the Client interface
type HashcatClient struct {
	config        *Config
	cmd           *exec.Cmd    // Legacy support
	activeSession CrackSession // Current active session
}

// NewClient creates a new hashcat client with the provided options
func NewClient(opts ...Option) (*HashcatClient, error) {
	config := DefaultConfig()

	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Verify hashcat binary exists and is executable
	if _, err := exec.LookPath(config.BinaryPath); err != nil {
		return nil, fmt.Errorf("hashcat binary not found or not executable: %w", err)
	}

	return &HashcatClient{
		config: config,
	}, nil
}

// GetDevices returns information about available devices
func (c *HashcatClient) GetDevices(ctx context.Context) (*models.DeviceList, error) {
	args := []string{"--backend-info", "--quiet"}

	output, err := c.executeCommand(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get device information: %w", err)
	}

	// Parse the device output
	devices, err := ParseDeviceOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device information: %w", err)
	}

	return devices, nil
}

// GetSupportedHashes returns information about supported hash types
func (c *HashcatClient) GetSupportedHashes(ctx context.Context) (*models.HashcatSupportedHashes, error) {
	args := []string{"--machine-readable", "--hash-info", "--quiet"}

	output, err := c.executeCommand(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported hash types: %w", err)
	}

	// Parse the output to extract hash types
	hashTypes, err := ParseHashInfoOutput(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hash types: %w", err)
	}
	return hashTypes, nil
}

// Benchmark performs a benchmark for the given hash type
func (c *HashcatClient) Benchmark(ctx context.Context, hashType int) (*models.HashcatBenchmarkResponse, error) {
	args := []string{
		"--hash-type", fmt.Sprintf("%d", hashType),
		"--benchmark",
		"--quiet",
	}

	output, err := c.executeCommand(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run benchmark: %w", err)
	}

	// Parse benchmark output
	return parseBenchmarkOutput(output, hashType), nil
}

// Crack attempts to crack the provided hash using the specified attack mode and options
func (c *HashcatClient) Crack(ctx context.Context, hash string, hashType int, attackMode int, mask string) (<-chan *models.Progress, error) {
	// Create options for the crack session
	options := &CrackOptions{
		HashType:        hashType,
		AttackMode:      attackMode,
		Mask:            mask,
		OptimizedKernel: true,
	}

	// Create a new crack session
	session, err := c.NewCrackSession(ctx, hash, options)
	if err != nil {
		return nil, err
	}
	c.activeSession = session

	// Start the cracking process
	if err := session.Start(); err != nil {
		return nil, err
	}

	// Return the progress channel
	return session.Progress(), nil
}

// CrackFile attempts to crack hashes in the specified file
func (c *HashcatClient) CrackFile(ctx context.Context, hashFile *models.HashFile, attackMode int, mask string) (<-chan *models.Progress, error) {
	// Create options for the crack session
	options := &CrackOptions{
		HashType:        hashFile.HashType,
		AttackMode:      attackMode,
		Mask:            mask,
		OptimizedKernel: true,
	}

	// Create a new crack file session
	session, err := c.NewCrackFileSession(ctx, hashFile.Path, options)
	if err != nil {
		return nil, err
	}
	c.activeSession = session

	// Start the cracking process
	if err := session.Start(); err != nil {
		return nil, err
	}

	// Return the progress channel
	return session.Progress(), nil
}

// Stop interrupts a running cracking session
func (c *HashcatClient) Stop(ctx context.Context) error {
	// First try to stop an active session if one exists
	if c.activeSession != nil {
		err := c.activeSession.Stop()
		c.activeSession = nil
		return err
	}

	// Fall back to legacy method for backwards compatibility
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}

	return nil
}

// executeCommand is a helper method to execute hashcat commands
func (c *HashcatClient) executeCommand(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c.config.BinaryPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}

func parseBenchmarkOutput(output string, hashType int) *models.HashcatBenchmarkResponse {
	scanner := bufio.NewScanner(strings.NewReader(output))

	response := &models.HashcatBenchmarkResponse{
		Benchmarks: []models.Benchmark{},
		Summary:    models.BenchmarkSummary{},
	}

	hashModeFound := false

	// Regular expressions for parsing
	hashModeRegex := regexp.MustCompile(`\* Hash-Mode (\d+) \((.+)\)`)
	speedRegex := regexp.MustCompile(`Speed\.#(\d+)\.+:\s+([0-9.]+)\s+([GMk]H/s).*?@\s+Accel:(\d+)\s+Loops:(\d+)\s+Thr:(\d+)\s+Vec:(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line contains hash mode information
		if matches := hashModeRegex.FindStringSubmatch(line); matches != nil {
			hashMode, _ := strconv.Atoi(matches[1])
			hashName := matches[2]

			hashModeFound = true
			benchmark := models.Benchmark{
				HashMode:      hashMode,
				HashName:      hashName,
				DeviceResults: []models.BenchmarkResult{},
			}
			response.Benchmarks = append(response.Benchmarks, benchmark)
		}

		// Check if line contains speed information
		if matches := speedRegex.FindStringSubmatch(line); matches != nil && hashModeFound {
			deviceID, _ := strconv.Atoi(matches[1])
			speed, _ := strconv.ParseFloat(matches[2], 64)
			speedUnit := matches[3]
			accel, _ := strconv.Atoi(matches[4])
			loops, _ := strconv.Atoi(matches[5])
			threads, _ := strconv.Atoi(matches[6])
			vecSize, _ := strconv.Atoi(matches[7])

			// Extract time per hash value (ms) if available
			timePerHash := 0.0
			timePerHashRegex := regexp.MustCompile(`\(([\d.]+)ms\)`)
			if timeMatches := timePerHashRegex.FindStringSubmatch(line); timeMatches != nil {
				timePerHash, _ = strconv.ParseFloat(timeMatches[1], 64)
			}

			result := models.BenchmarkResult{
				DeviceID:     deviceID,
				Speed:        speed,
				SpeedUnit:    speedUnit,
				TimePerHash:  timePerHash,
				Acceleration: accel,
				Loops:        loops,
				Threads:      threads,
				VectorSize:   vecSize,
			}

			// Add result to the last benchmark
			if len(response.Benchmarks) > 0 {
				response.Benchmarks[len(response.Benchmarks)-1].DeviceResults = append(response.Benchmarks[len(response.Benchmarks)-1].DeviceResults, result)
			}
		}
	}
	return response
}
