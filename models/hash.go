package models

// HashType represents a hashcat hash type
type HashType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// HashMode represents a hashcat hash mode with its settings
type HashMode struct {
	HashType
	IsOptimized bool `json:"is_optimized"`
	IsSalted    bool `json:"is_salted"`
}

// HashFile represents a file containing hashes to be cracked
type HashFile struct {
	Path     string `json:"path"`
	HashType int    `json:"hash_type"`
	Count    int    `json:"count,omitempty"`
}

// HashcatSupportedHashes represents the response from hashcat when listing supported hash types
type HashcatSupportedHashes struct {
	HashTypes []HashType `json:"hash_types"`
}
