package custodian

import (
	"context"
	"io/fs"
)

type DependencyTree interface {
	GetModule(moduleIdentifier string) (Module, bool)
	GenerateLockFile() []byte
	RootIdentifier() string
}
type Module interface {
	GetDependencyModule(dependencyName string, dt DependencyTree) (Module, string)
	DependencyList() []string
	FileSystem() fs.FS
	Identifier() string
}
type Resolver interface {
	Resolve(ctx context.Context, moduleIdentifier string) (Module, error)
}
