package modules

import (
	"fmt"
	"io"
	"path"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"

	"github.com/google/go-jsonnet"
)

type Transformer func(foundAt string, data []byte) ([]byte, error)

type GitImporter struct {
	Transformers   []Transformer
	DependencyTree custodian.DependencyTree
}

func (importer *GitImporter) Import(importedFrom string, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	var detectedModule custodian.Module
	var detectedPath string

	// Parse the importedFrom to get the module identifier
	sourceModuleIdentifier, sourceFilePath := utils.ParseImportedFrom(importedFrom)
	if sourceModuleIdentifier == "" {
		sourceModuleIdentifier = importer.DependencyTree.RootIdentifier()
	}
	if sourceFilePath == "" {
		sourceFilePath = DefaultRootFile
	}

	// Load the source module from the dependency tree
	sourceModule, exists := importer.DependencyTree.GetModule(sourceModuleIdentifier)
	if !exists {
		return jsonnet.MakeContents(""), "", fmt.Errorf("module not found: %s", sourceModuleIdentifier)
	}

	// Parse the imported path to get the dependency name and path
	dependencyName, dependencyPath := utils.ParseImportedPath(importedPath)

	dependencyModule, dependencyModuleIdentificer := sourceModule.GetDependencyModule(dependencyName, importer.DependencyTree)
	// If the dependency module is not found, it means it's a local import
	if dependencyModuleIdentificer == "" {
		detectedModule = sourceModule
		// handle relative imports
		if utils.IsRelativeImport(importedPath) {
			detectedPath = path.Join(path.Dir(sourceFilePath), importedPath)
		} else {
			detectedPath = importedPath
		}
		foundAt = utils.BuildFoundAtPath(sourceModuleIdentifier, detectedPath)
	} else if dependencyModule != nil {
		// import from a dependency module found in the dependency tree
		detectedModule = dependencyModule
		detectedPath = dependencyPath
		foundAt = utils.BuildFoundAtPath(dependencyModuleIdentificer, detectedPath)
	} else {
		// Dependency module not found in the dependency tree
		return jsonnet.MakeContents(""), "", fmt.Errorf("dependency module not found: %s", dependencyModuleIdentificer)
	}

	file, err := detectedModule.FileSystem().Open(detectedPath)
	if err != nil {
		return jsonnet.MakeContents(""), "", err
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		return jsonnet.MakeContents(""), "", err
	}

	// Apply transformers
	for _, transformer := range importer.Transformers {
		fileData, err = transformer(foundAt, fileData)
		if err != nil {
			return jsonnet.MakeContents(""), "", err
		}
	}

	// Return the contents
	contents = jsonnet.MakeContentsRaw(fileData)
	return contents, foundAt, nil
}

func (importer *GitImporter) AddTransformer(transformer Transformer) {
	importer.Transformers = append(importer.Transformers, transformer)
}
