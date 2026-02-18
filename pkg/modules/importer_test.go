package modules

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"io/fs"
	"path"
	"testing"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"

	"github.com/google/go-jsonnet"
)

type TestModuleResolver struct {
}

//go:embed "fixtures/repoA@v0.0.1/sub1/*" "fixtures/repoA@v0.0.1/*" "fixtures/repoB@v0.0.2/*"
var fixtures embed.FS

func (c *TestModuleResolver) Resolve(ctx context.Context, moduleIdentifier string) (custodian.Module, error) {
	rFs, err := fs.Sub(fixtures, path.Join("fixtures", moduleIdentifier))
	if err != nil {
		return nil, err
	}
	return NewModuleFromFS(moduleIdentifier, rFs)
}

func TestGitImporter_Import(t *testing.T) {

	resolver := &TestModuleResolver{}

	root, err := resolver.Resolve(context.Background(), "repoA@v0.0.1")
	if err != nil {
		t.Fatalf("Failed to load root module: %v", err)
	}

	dependencyTree, err := NewDependencyTree(root, resolver)
	if err != nil {
		t.Fatalf("Failed to build dependency tree: %v", err)
	}

	importer := GitImporter{
		DependencyTree: dependencyTree,
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		importedFrom string
		importedPath string
		want         jsonnet.Contents
		want2        string
		wantErr      bool
	}{
		{
			name:         "import from local module",
			importedFrom: "",
			importedPath: "./sub1/sub1.libsonnet",
			want:         jsonnet.MakeContents("{\n    a: 1,\n    b: 2,\n    c: 3,\n}"),
			want2:        "repoA@v0.0.1" + utils.ModuleIdentifierSeparator + "sub1/sub1.libsonnet",
			wantErr:      false,
		},
		{
			name:         "import from dependency module",
			importedFrom: "repoA@v0.0.1" + utils.ModuleIdentifierSeparator + "main.jsonnet",
			importedPath: "modB/main.jsonnet",
			want:         jsonnet.MakeContents("{\n    f: 5,\n    g: 6,\n    h: 7,\n}"),
			want2:        "repoB@v0.0.2" + utils.ModuleIdentifierSeparator + "main.jsonnet",
			wantErr:      false,
		},
		{
			name:         "import from local module relative path",
			importedFrom: "repoA@v0.0.1" + utils.ModuleIdentifierSeparator + "sub1/xalala.libsonnet",
			importedPath: "./sub1.libsonnet",
			want:         jsonnet.MakeContents("{\n    a: 1,\n    b: 2,\n    c: 3,\n}"),
			want2:        "repoA@v0.0.1" + utils.ModuleIdentifierSeparator + "sub1/sub1.libsonnet",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotErr := importer.Import(tt.importedFrom, tt.importedPath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Import() failed: %v, test case: %s", gotErr, tt.name)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Import() succeeded unexpectedly")
			}

			if !bytes.Equal(got.Data(), tt.want.Data()) {
				t.Errorf("Import() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("Import() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
