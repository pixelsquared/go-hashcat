package hashcat

import (
	"os/exec"
)

// Config stores the configuration for the hashcat client
type Config struct {
	// BinaryPath is the path to the hashcat executable
	BinaryPath string

	// OutputDir is the directory where hashcat will store output files
	OutputDir string

	// DefaultAttackMode is the default attack mode to use
	DefaultAttackMode int

	// DefaultHashType is the default hash type to use
	DefaultHashType int

	// AdditionalOptions contains additional command-line options to pass to hashcat
	AdditionalOptions []string
}

// Option is a function that configures a Config
type Option func(*Config) error

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	// Find hashcat in PATH by default
	hashcatPath, _ := exec.LookPath("hashcat")
	if hashcatPath == "" {
		hashcatPath = "hashcat" // Fallback to just the name and assume it's in PATH
	}

	return &Config{
		BinaryPath:        hashcatPath,
		OutputDir:         "./hashcat-output",
		DefaultAttackMode: 0, // Straight mode
		DefaultHashType:   0, // MD5
	}
}

// WithBinaryPath sets the path to the hashcat executable
func WithBinaryPath(path string) Option {
	return func(c *Config) error {
		if path == "" {
			return ErrInvalidBinaryPath
		}

		c.BinaryPath = path
		return nil
	}
}

// WithOutputDir sets the directory where hashcat will store output files
func WithOutputDir(dir string) Option {
	return func(c *Config) error {
		if dir == "" {
			return ErrInvalidOutputDir
		}

		c.OutputDir = dir
		return nil
	}
}

// WithDefaultAttackMode sets the default attack mode
func WithDefaultAttackMode(mode int) Option {
	return func(c *Config) error {
		if mode < 0 || mode > 9 {
			return ErrInvalidAttackMode
		}

		c.DefaultAttackMode = mode
		return nil
	}
}

// WithDefaultHashType sets the default hash type
func WithDefaultHashType(hashType int) Option {
	return func(c *Config) error {
		if hashType < 0 {
			return ErrInvalidHashType
		}

		c.DefaultHashType = hashType
		return nil
	}
}

// WithAdditionalOptions adds extra command-line options to be passed to hashcat
func WithAdditionalOptions(options ...string) Option {
	return func(c *Config) error {
		c.AdditionalOptions = append(c.AdditionalOptions, options...)
		return nil
	}
}
