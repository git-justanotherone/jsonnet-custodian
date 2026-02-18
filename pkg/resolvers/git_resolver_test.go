package resolvers

import (
	"os"
	"testing"
)

func Test_gitResolver_getModule(t *testing.T) {

	dir, err := os.MkdirTemp("", "gittest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	f, err := NewGitResolver(dir)
	if err != nil {
		t.Fatalf("Failed to create GitResolver: %v", err)
	}
	gf := f.(*gitResolver)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		moduleIdentifier string
		want             string
		wantErr          bool
	}{
		{
			name:             "Test cloning a valid git module with tag version",
			moduleIdentifier: "github.com/google/go-jsonnet@v0.21.0",
			want:             "github.com/google/go-jsonnet@v0.21.0",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module with commit hash matching a tag",
			moduleIdentifier: "github.com/google/go-jsonnet@67968688d9952f506",
			want:             "github.com/google/go-jsonnet@v0.21.0",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module from HEAD",
			moduleIdentifier: "github.com/getsops/gopgagent",
			want:             "github.com/getsops/gopgagent@v0.0.0-20241224165529-7044f28e491e",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module with commit hash not matching a tag",
			moduleIdentifier: "github.com/go-git/gcfg@0495f244e4712b22418b489b2079d3a1598f434f",
			want:             "github.com/go-git/gcfg@v2.0.2-0.20250606171425-0495f244e471",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module with a long pseudo-version",
			moduleIdentifier: "github.com/go-git/gcfg@v2.0.1-0.20240810071503-633fd26b3faaa8e7a2548e4b4b0b1a07618eba86",
			want:             "github.com/go-git/gcfg@v2.0.1-0.20240810071503-633fd26b3faa",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module with a valid pseudo-version",
			moduleIdentifier: "github.com/go-git/gcfg@v2.0.1-0.20240810071503-633fd26b3faa",
			want:             "github.com/go-git/gcfg@v2.0.1-0.20240810071503-633fd26b3faa",
			wantErr:          false,
		},
		{
			name:             "Test cloning a valid git module with a valid pseudo-version",
			moduleIdentifier: "github.com/go-git/gcfg@v2.0.1-0.20240810071503-633as26b3faa",
			want:             "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			got, gotErr := gf.getModule(tt.moduleIdentifier)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getModule() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getModule() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("getModule() = %v, want %v", got, tt.want)
			}
		})
	}
	// t.Error(dir)
}
