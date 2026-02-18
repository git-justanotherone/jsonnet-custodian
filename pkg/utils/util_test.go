package utils

import (
	"testing"
)

func TestIsRelativeImport(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		importedPath string
		want         bool
	}{
		{
			name:         "relative import with ./",
			importedPath: "./some/path/file.jsonnet",
			want:         true,
		},
		{
			name:         "relative import with ../",
			importedPath: "../some/path/file.jsonnet",
			want:         true,
		},
		{
			name:         "absolute import",
			importedPath: "some/path/file.jsonnet",
			want:         false,
		},
		{
			name:         "absolute import with prefix slash",
			importedPath: "/some/path/file.jsonnet",
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRelativeImport(tt.importedPath)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("IsRelativeImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseImportedFrom(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		importedFrom string
		want         string
		want2        string
	}{
		{
			name:         "basic test",
			importedFrom: "repoA@v1.0.0" + ModuleIdentifierSeparator + "some/path/file.jsonnet",
			want:         "repoA@v1.0.0",
			want2:        "some/path/file.jsonnet",
		},
		{
			name:         "no path test",
			importedFrom: "repoB@v2.0.0",
			want:         "repoB@v2.0.0",
			want2:        "",
		},
		{
			name:         "empty import test",
			importedFrom: "",
			want:         "",
			want2:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := ParseImportedFrom(tt.importedFrom)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ParseImportedFrom() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseImportedFrom() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestParseImportedPath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		importedPath string
		want         string
		want2        string
	}{
		{
			name:         "basic test",
			importedPath: "repoB/other/path/otherfile.jsonnet",
			want:         "repoB",
			want2:        "other/path/otherfile.jsonnet",
		},
		{
			name:         "no path test",
			importedPath: "repoC",
			want:         "repoC",
			want2:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := ParseImportedPath(tt.importedPath)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ParseImportedPath() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseImportedPath() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
