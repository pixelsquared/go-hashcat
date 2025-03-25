package hashcat

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pixelsquared/go-hashcat/models"
)

// ParseHashInfoOutput parses the JSON output from hashcat --hash-info command
// This is exported so it can be used in client.go
func ParseHashInfoOutput(output string) (*models.HashcatSupportedHashes, error) {
	// Parse the raw JSON output
	var rawData map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(output), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse hash info JSON: %w", err)
	}

	// Convert the raw data to HashType structs
	hashTypes := make([]models.HashType, 0, len(rawData))

	for hashIDStr, hashData := range rawData {
		// Convert hash ID from string to int
		hashID, err := strconv.Atoi(hashIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid hash ID format '%s': %w", hashIDStr, err)
		}

		// Extract the hash type information
		name, _ := hashData["name"].(string)
		category, _ := hashData["category"].(string)

		// Create HashType struct
		hashType := models.HashType{
			ID:          hashID,
			Name:        name,
			Category:    category,
			Description: "", // Description field is not provided in the hash-info output
		}

		hashTypes = append(hashTypes, hashType)
	}

	return &models.HashcatSupportedHashes{
		HashTypes: hashTypes,
	}, nil
}

// FindHashTypeByID returns a HashType with the given ID or nil if not found
func FindHashTypeByID(hashes *models.HashcatSupportedHashes, id int) *models.HashType {
	for _, hashType := range hashes.HashTypes {
		if hashType.ID == id {
			return &hashType
		}
	}
	return nil
}

// FindHashTypeByName returns a HashType with the given name or nil if not found
// Note: This does a case-sensitive exact match
func FindHashTypeByName(hashes *models.HashcatSupportedHashes, name string) *models.HashType {
	for _, hashType := range hashes.HashTypes {
		if hashType.Name == name {
			return &hashType
		}
	}
	return nil
}

// GetHashTypeIDByName returns the ID of a hash type with the given name or -1 if not found
func GetHashTypeIDByName(hashes *models.HashcatSupportedHashes, name string) int {
	hashType := FindHashTypeByName(hashes, name)
	if hashType == nil {
		return -1
	}
	return hashType.ID
}

// FindHashTypesByCategory returns all hash types in the given category
func FindHashTypesByCategory(hashes *models.HashcatSupportedHashes, category string) []models.HashType {
	var result []models.HashType
	for _, hashType := range hashes.HashTypes {
		if strings.EqualFold(hashType.Category, category) {
			result = append(result, hashType)
		}
	}
	return result
}
