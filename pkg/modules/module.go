package modules

import (
	"encoding/json"
	"io/fs"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
)

const (
	ModuleFileName = "custodian.json"
)

type ModuleFile struct {
	Module  string            `json:"module"`
	Require map[string]string `json:"require"`
}

type module struct {
	dependencies map[string]string // local name -> module_identifier
	fileSystem   fs.FS
	identifier   string
}

func (m *module) GetDependencyModule(dependencyName string, dt custodian.DependencyTree) (custodian.Module, string) {
	moduleIdentifier, exists := m.dependencies[dependencyName]
	if !exists {
		return nil, ""
	}

	dModule, exists := dt.GetModule(moduleIdentifier)
	if !exists {
		return nil, moduleIdentifier
	}

	return dModule, moduleIdentifier
}

func (m *module) DependencyList() []string {
	depList := make([]string, 0, len(m.dependencies))
	for _, depId := range m.dependencies {
		depList = append(depList, depId)
	}
	return depList
}

func (m *module) FileSystem() fs.FS {
	return m.fileSystem
}

func (m *module) Identifier() string {
	return m.identifier
}

func NewModuleFromFS(moduleIdentifier string, moduleFS fs.FS) (custodian.Module, error) {

	moduleFile, err := moduleFS.Open(ModuleFileName)
	if err != nil {
		// If there is no module file, return an empty module
		return &module{
			dependencies: map[string]string{},
			fileSystem:   moduleFS,
			identifier:   moduleIdentifier,
		}, nil
	}
	defer moduleFile.Close()

	moduleData, err := ParseModuleFile(moduleFile)
	if err != nil {
		return nil, err
	}

	return &module{
		dependencies: moduleData.Require,
		fileSystem:   moduleFS,
		identifier:   moduleIdentifier,
	}, nil

}

func ParseModuleFile(moduleFile fs.File) (*ModuleFile, error) {
	moduleData := &ModuleFile{}
	if err := json.NewDecoder(moduleFile).Decode(moduleData); err != nil {
		return nil, err
	}
	return moduleData, nil
}

func SerializeModuleFile(moduleData *ModuleFile) ([]byte, error) {
	data, err := json.MarshalIndent(moduleData, "", "    ")
	if err != nil {
		return nil, err
	}
	return data, nil
}
