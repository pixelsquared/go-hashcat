package hashcat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/pixelsquared/go-hashcat/models"
)

// CrackSession represents an active hashcat cracking session
type CrackSession interface {
	// Start initiates the cracking process
	Start() error

	// Stop terminates the cracking process
	Stop() error

	// Progress returns a channel that receives progress updates
	Progress() <-chan *models.Progress

	// Wait blocks until the cracking process completes
	Wait() error

	// Results returns the final cracking results
	Results() ([]*models.CrackedHash, error)
}

// HashcatCrackSession implements the CrackSession interface
type HashcatCrackSession struct {
	client       *HashcatClient
	cmd          *exec.Cmd
	progressChan chan *models.Progress
	resultsChan  chan *models.CrackedHash
	ctx          context.Context
	cancel       context.CancelFunc
	mutex        sync.Mutex
	isRunning    bool
	hashFile     string
	outputFile   string
	potFile      string
	sessionName  string
	results      []*models.CrackedHash
	errorChan    chan error
	finalError   error
	wg           sync.WaitGroup
}

// CrackOptions defines parameters for a cracking session
type CrackOptions struct {
	HashType        int      // Hash type ID
	AttackMode      int      // Attack mode (0=dict, 1=combi, 3=mask, etc.)
	Mask            string   // Mask for mask attack or wordlist for dictionary attack
	Rules           []string // Rules to apply
	OptimizedKernel bool     // Use optimized kernels if available (default: true)
	Workload        int      // Workload profile (1=low, 2=default, 3=high, 4=nightmare)
	DeviceIDs       []int    // Specific device IDs to use (empty=all devices)
}

// NewCrackSession creates a new CrackSession for cracking a single hash
func (c *HashcatClient) NewCrackSession(ctx context.Context, hash string, options *CrackOptions) (CrackSession, error) {
	// Create temporary files for this session
	sessionName := fmt.Sprintf("hashcat-%d", time.Now().UnixNano())
	hashFile, err := createTempFile(sessionName+"-hash.txt", hash+"\n")
	if err != nil {
		return nil, fmt.Errorf("failed to create hash file: %w", err)
	}

	outputFile := hashFile + ".out"
	potFile := hashFile + ".pot"

	// Create derived context with cancellation and options
	optionsCtx := context.WithValue(ctx, "options", options)
	ctx, cancel := context.WithCancel(optionsCtx)

	return &HashcatCrackSession{
		client:       c,
		progressChan: make(chan *models.Progress, 10),
		resultsChan:  make(chan *models.CrackedHash, 100),
		ctx:          ctx,
		cancel:       cancel,
		hashFile:     hashFile,
		outputFile:   outputFile,
		potFile:      potFile,
		sessionName:  sessionName,
		results:      []*models.CrackedHash{},
		errorChan:    make(chan error, 1),
	}, nil
}

// NewCrackFileSession creates a new CrackSession for cracking multiple hashes from a file
func (c *HashcatClient) NewCrackFileSession(ctx context.Context, hashFilePath string, options *CrackOptions) (CrackSession, error) {
	// Verify hash file exists
	if _, err := os.Stat(hashFilePath); err != nil {
		return nil, fmt.Errorf("hash file not found: %w", err)
	}

	sessionName := fmt.Sprintf("hashcat-%d", time.Now().UnixNano())
	outputFile := hashFilePath + ".out"
	potFile := hashFilePath + ".pot"

	// Create derived context with cancellation and options
	optionsCtx := context.WithValue(ctx, "options", options)
	ctx, cancel := context.WithCancel(optionsCtx)

	return &HashcatCrackSession{
		client:       c,
		progressChan: make(chan *models.Progress, 10),
		resultsChan:  make(chan *models.CrackedHash, 100),
		ctx:          ctx,
		cancel:       cancel,
		hashFile:     hashFilePath,
		outputFile:   outputFile,
		potFile:      potFile,
		sessionName:  sessionName,
		results:      []*models.CrackedHash{},
		errorChan:    make(chan error, 1),
	}, nil
}

// Start initiates the cracking process
func (s *HashcatCrackSession) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return fmt.Errorf("cracking session already running")
	}

	// Extract options from context
	options, ok := s.ctx.Value("options").(*CrackOptions)
	if !ok {
		options = &CrackOptions{
			HashType:        0,
			AttackMode:      0,
			OptimizedKernel: true,
		}
	}

	// Construct command arguments
	args := []string{
		fmt.Sprintf("--hash-type=%d", options.HashType),
		fmt.Sprintf("--attack-mode=%d", options.AttackMode),
		"--quiet",
		"--status",
		"--status-json",
		"--status-timer", "1",
		"--session", s.sessionName,
		"--outfile", s.outputFile,
		"--potfile-path", s.potFile,
		s.hashFile,
		"", // Mask or wordlist, will be properly set when using options
	}

	// Add optimized kernel if requested
	if options.OptimizedKernel {
		args = append([]string{"--optimized-kernel-enable"}, args...)
	}

	// Set workload profile if specified
	if options.Workload > 0 {
		args = append(args[:len(args)-1], fmt.Sprintf("--workload-profile=%d", options.Workload))
	}

	// Add rules if specified
	for _, rule := range options.Rules {
		args = append(args[:len(args)-1], fmt.Sprintf("--rule=%s", rule))
	}

	// Add device IDs if specified
	for _, deviceID := range options.DeviceIDs {
		args = append(args[:len(args)-1], fmt.Sprintf("--opencl-device-types=%d", deviceID))
	}

	// Replace the last empty arg with the mask/wordlist
	args[len(args)-1] = options.Mask

	// Set up and execute command
	cmd := exec.CommandContext(s.ctx, s.client.config.BinaryPath, args...)
	s.cmd = cmd

	// Get stdout pipe for progress updates
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	// Get stderr pipe for error messages
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start hashcat: %w", err)
	}

	s.isRunning = true

	// Process stdout for progress updates
	s.wg.Add(1)
	go s.processOutput(stdout, stderr)

	// Process results from output file
	s.wg.Add(1)
	go s.monitorResults()

	return nil
}

// processOutput reads and parses the JSON output from hashcat
func (s *HashcatCrackSession) processOutput(stdout, stderr io.ReadCloser) {
	defer s.wg.Done()
	defer close(s.progressChan)

	// Create scanner for stdout
	scanner := bufio.NewScanner(stdout)

	// Create error scanner for stderr
	errScanner := bufio.NewScanner(stderr)

	// Start a goroutine to collect stderr output
	go func() {
		var errOutput string
		for errScanner.Scan() {
			errOutput += errScanner.Text() + "\n"
		}

		if errOutput != "" {
			select {
			case s.errorChan <- fmt.Errorf("hashcat error: %s", errOutput):
			default:
			}
		}
	}()

	// Process each line from stdout
	for scanner.Scan() {
		line := scanner.Text()

		// Parse JSON progress update
		var progress models.Progress
		if err := json.Unmarshal([]byte(line), &progress); err == nil {
			// Send progress update through channel
			select {
			case s.progressChan <- &progress:
			default:
				// Channel buffer is full, discard oldest value and try again
				<-s.progressChan
				s.progressChan <- &progress
			}

			// If status indicates completion, break
			if progress.Status == models.StatusCracked ||
				progress.Status == models.StatusExhausted ||
				progress.Status == models.StatusAborted {
				break
			}
		}
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		select {
		case s.errorChan <- fmt.Errorf("error reading output: %w", err):
		default:
		}
	}
}

// monitorResults monitors the output file for cracked hashes
func (s *HashcatCrackSession) monitorResults() {
	defer s.wg.Done()
	defer close(s.resultsChan)

	// Check output file for results periodically
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastSize int64

	for {
		select {
		case <-ticker.C:
			// Get file info
			info, err := os.Stat(s.outputFile)
			if err != nil {
				if !os.IsNotExist(err) {
					// Only report error if it's not just that the file doesn't exist yet
					select {
					case s.errorChan <- fmt.Errorf("error checking output file: %w", err):
					default:
					}
				}
				continue
			}

			// If file size hasn't changed, skip
			if info.Size() <= lastSize {
				continue
			}

			// Open file and seek to last read position
			file, err := os.Open(s.outputFile)
			if err != nil {
				select {
				case s.errorChan <- fmt.Errorf("error opening output file: %w", err):
				default:
				}
				continue
			}

			if _, err := file.Seek(lastSize, 0); err != nil {
				file.Close()
				select {
				case s.errorChan <- fmt.Errorf("error seeking in output file: %w", err):
				default:
				}
				continue
			}

			// Read new content
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// Parse cracked hash
				parts := splitHashResult(line)
				if len(parts) >= 2 {
					hash := parts[0]
					password := parts[1]

					result := &models.CrackedHash{
						Hash:     hash,
						Password: password,
						Time:     time.Now().Unix(),
					}

					// Send through channel and add to results slice
					s.resultsChan <- result

					s.mutex.Lock()
					s.results = append(s.results, result)
					s.mutex.Unlock()
				}
			}

			lastSize = info.Size()
			file.Close()

		case <-s.ctx.Done():
			return
		}
	}
}

// Progress returns the progress channel
func (s *HashcatCrackSession) Progress() <-chan *models.Progress {
	return s.progressChan
}

// Stop terminates the cracking process
func (s *HashcatCrackSession) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	defer s.cleanup() // Clean up temporary files

	if !s.isRunning {
		return nil
	}

	// Cancel context to stop all operations
	s.cancel()

	// Kill the process if it's still running
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}

	return nil
}

// Wait blocks until the cracking process completes
func (s *HashcatCrackSession) Wait() error {
	// Wait for all goroutines to finish
	s.wg.Wait()
	defer s.cleanup() // Clean up temporary files

	// Check for errors
	select {
	case err := <-s.errorChan:
		s.finalError = err
	default:
		// No error occurred
	}

	return s.finalError
}

// Results returns the cracked hashes
func (s *HashcatCrackSession) Results() ([]*models.CrackedHash, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.results, s.finalError
}

// cleanup removes temporary files created for the session
func (s *HashcatCrackSession) cleanup() {
	// Only remove temp hash file if we created it
	// Don't try to clean up if we haven't set a hash file yet
	if s.hashFile == "" || strings.HasPrefix(s.hashFile, os.TempDir()) {
		os.Remove(s.hashFile)
	}

	// Remove output and pot files
	os.Remove(s.outputFile)
	os.Remove(s.potFile)
}

// Helper function to create a temporary file with content
func createTempFile(name, content string) (string, error) {
	file, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}

// Helper function to split hash:password results
func splitHashResult(line string) []string {
	var parts []string
	var current string
	inEscape := false

	for _, c := range line {
		if c == '\\' && !inEscape {
			inEscape = true
			continue
		}

		if c == ':' && !inEscape {
			parts = append(parts, current)
			current = ""
			continue
		}

		current += string(c)
		inEscape = false
	}

	parts = append(parts, current)
	return parts
}
