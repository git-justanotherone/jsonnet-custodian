package modules

import (
	"context"
	"encoding/json"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"
)

const (
	DefaultRootFile = "main.jsonnet"
)

type dependencyTree struct {
	rootIdentifier string
	modules        map[string]custodian.Module // module_identifier -> Module
}

func (dt *dependencyTree) GetModule(moduleIdentifier string) (custodian.Module, bool) {
	module, exists := dt.modules[moduleIdentifier]
	return module, exists
}

func (dt *dependencyTree) RootIdentifier() string {
	return dt.rootIdentifier
}

func (dt *dependencyTree) GenerateLockFile() []byte {
	lockData := make([]string, 0, len(dt.modules))
	for _, module := range dt.modules {
		// Only include non-local modules in the lock file
		if !utils.IsLocalPath(module.Identifier()) {
			lockData = append(lockData, module.Identifier())
		}
	}
	lockFileBytes, _ := json.MarshalIndent(lockData, "", "  ")
	return lockFileBytes
}

func NewDependencyTree(root custodian.Module, resolver custodian.Resolver) (custodian.DependencyTree, error) {
	modules := make(map[string]custodian.Module)
	modules[root.Identifier()] = root

	hasNewModule := true
	for hasNewModule {
		hasNewModule = false

		for _, module := range modules {
			for _, depModuleId := range module.DependencyList() {
				// If we have not seen this module yet, fetch it and mark that we have new modules to process
				if _, exists := modules[depModuleId]; !exists {
					hasNewModule = true
					depModule, err := resolver.Resolve(context.Background(), depModuleId)
					if err != nil {
						return nil, err
					}
					modules[depModuleId] = depModule
				}
			}
		}
	}

	return &dependencyTree{modules: modules, rootIdentifier: root.Identifier()}, nil
}
