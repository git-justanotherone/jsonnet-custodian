package utils

import (
	"os"
	"strings"
)

const (
	ModuleIdentifierSeparator = ":mod-sep:"
	ENV_FILE_SUFFIX           = "_FILE"
)

func IsRelativeImport(importedPath string) bool {
	return strings.HasPrefix(importedPath, "./") || strings.HasPrefix(importedPath, "../")
}

func ParseImportedFrom(importedFrom string) (moduleIdentifier, path string) {
	if len(importedFrom) == 0 {
		return "", ""
	}
	split := strings.SplitN(importedFrom, ModuleIdentifierSeparator, 2)
	if len(split) < 2 {
		return split[0], ""
	}
	return split[0], split[1]
}

func BuildFoundAtPath(moduleIdentifier, filePath string) string {
	return moduleIdentifier + ModuleIdentifierSeparator + filePath
}

func ParseImportedPath(importedPath string) (repository string, filepath string) {
	split := strings.SplitN(importedPath, "/", 2)
	if len(split) < 2 {
		return split[0], ""
	}
	return split[0], split[1]
}

func GetEnv(envVar, defaultValue string) string {
	if filePath := os.Getenv(envVar + ENV_FILE_SUFFIX); filePath != "" {
		if content, err := os.ReadFile(filePath); err == nil {
			return string(content)
		}
	}

	if envValue, exists := os.LookupEnv(envVar); exists {
		return envValue
	}

	return defaultValue
}

func GetEnvOrEmpty(envVar string) string {
	return GetEnv(envVar, "")
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsLocalPath(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") || strings.HasPrefix(path, "/") || path == "." || path == ".."
}
