package transformers

import (
	"strings"

	"github.com/getsops/sops/v3/decrypt"
)

// getSopsFormat determines if the filename indicates a sops-encrypted file
// and returns the corresponding format string and a boolean indicating if it is sops-encrypted.
func getSopsFormat(filename string) (format string, isSops bool) {
	switch {
	case strings.HasSuffix(filename, ".sops.json"):
		return "json", true
	case strings.HasSuffix(filename, ".sops.yaml"), strings.HasSuffix(filename, ".sops.yml"):
		return "yaml", true
	case strings.HasSuffix(filename, ".sops.ini"):
		return "ini", true
	case strings.HasSuffix(filename, ".sops.env"), strings.HasSuffix(filename, ".sops.dotenv"):
		return "dotenv", true
	case strings.Contains(filename, ".sops."):
		return "binary", true
	default:
		return "", false
	}
}

// SopsDecryptorTransformer is a transformer function that decrypts sops-encrypted data
// based on the imported path's file extension.
func SopsDecryptorTransformer(foundAt string, data []byte) ([]byte, error) {
	sopsFormat, isSops := getSopsFormat(foundAt)
	if !isSops {
		return data, nil
	}

	return decrypt.Data(data, sopsFormat)
}
