package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/modules"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/resolvers"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/transformers"

	"github.com/google/go-jsonnet"
)

const (
	MODULE_CACHE_DIR = "/tmp/jnetx/modules"
	LOCK_FILE_NAME   = "module.lock"
)

func GetDependencyTree() (custodian.DependencyTree, error) {
	// Create a Resolver
	moduleResolver, err := resolvers.NewResolver(MODULE_CACHE_DIR)
	if err != nil {
		return nil, err
	}

	root, err := moduleResolver.Resolve(context.Background(), ".")
	if err != nil {
		return nil, err
	}
	dt, err := modules.NewDependencyTree(root, moduleResolver)
	if err != nil {
		return nil, err
	}

	return dt, nil
}

func GetModule(moduleIdentifier string) (string, error) {

	// Create a Resolver
	moduleResolver, err := resolvers.NewResolver(MODULE_CACHE_DIR)
	if err != nil {
		return "", err
	}

	// Resolve the module identifier to get the filesystem and the resolved identifier
	module, err := moduleResolver.Resolve(context.Background(), moduleIdentifier)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stdout, "Module '%s' added/updated successfully in '%s'.\n", module.Identifier(), modules.ModuleFileName)
	return module.Identifier(), nil

}

func ConfigureVMExtensions(vm *jsonnet.VM) error {
	// Set up the GitImporter with the dependency tree.
	dt, err := GetDependencyTree()
	if err != nil {
		return err
	}
	importer := &modules.GitImporter{
		DependencyTree: dt,
	}
	// Add SOPS decryption transformer.
	importer.AddTransformer(transformers.SopsDecryptorTransformer)
	vm.Importer(importer)
	return nil
}
