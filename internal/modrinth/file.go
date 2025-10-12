/*
Copyright © 2025 Vincent De Borger <hello@vincentdeborger.be>
*/
package modrinth

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

// Get the Sha512 hash of a file
func getHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new hash interfact
	hash := sha512.New()

	// Copy data from file into hash function
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	// Get resulting hash
	hashBytes := hash.Sum(nil)

	return fmt.Sprintf("%x", hashBytes), nil
}
