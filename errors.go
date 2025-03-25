package hashcat

import (
	"errors"
	"fmt"
)

// Custom error values for various scenarios
var (
	ErrInvalidBinaryPath = errors.New("invalid hashcat binary path")
	ErrInvalidOutputDir  = errors.New("invalid output directory")
	ErrInvalidAttackMode = errors.New("invalid attack mode, must be between 0 and 9")
	ErrInvalidHashType   = errors.New("invalid hash type")
	ErrBinaryNotFound    = errors.New("hashcat binary not found or not executable")
	ErrExecutionFailed   = errors.New("hashcat execution failed")
	ErrInvalidHash       = errors.New("invalid hash format")
	ErrInvalidHashFile   = errors.New("invalid hash file")
)

// HashcatError represents a specific hashcat error with context
type HashcatError struct {
	Operation string // The operation that failed
	Err       error  // The underlying error
	Output    string // Command output if available
}

// Error implements the error interface
func (e *HashcatError) Error() string {
	if e.Output != "" {
		return fmt.Sprintf("hashcat %s failed: %v - output: %s", e.Operation, e.Err, e.Output)
	}
	return fmt.Sprintf("hashcat %s failed: %v", e.Operation, e.Err)
}

// Unwrap implements the error unwrapping interface
func (e *HashcatError) Unwrap() error {
	return e.Err
}

// NewHashcatError creates a new HashcatError
func NewHashcatError(operation string, err error, output string) *HashcatError {
	return &HashcatError{
		Operation: operation,
		Err:       err,
		Output:    output,
	}
}
