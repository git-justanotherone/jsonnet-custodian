package resolvers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/modules"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"
)

type localResolver struct{}

func (f *localResolver) Resolve(ctx context.Context, moduleIdentifier string) (custodian.Module, error) {
	fPath, err := filepath.Abs(moduleIdentifier)
	if err != nil {
		return nil, err
	}
	rFs := os.DirFS(fPath)
	module, err := modules.NewModuleFromFS(moduleIdentifier, rFs)
	if err != nil {
		return nil, err
	}
	return module, nil
}

type chainResolver struct {
	localResolver custodian.Resolver
	gitResolver   custodian.Resolver
}

func (f *chainResolver) Resolve(ctx context.Context, moduleIdentifier string) (custodian.Module, error) {
	if utils.IsLocalPath(moduleIdentifier) {
		return f.localResolver.Resolve(ctx, moduleIdentifier)
	} else {
		return f.gitResolver.Resolve(ctx, moduleIdentifier)
	}
}

func NewResolver(targetDir string) (custodian.Resolver, error) {
	gitResolver, err := NewGitResolver(targetDir)
	if err != nil {
		return nil, err
	}

	return &chainResolver{
		localResolver: &localResolver{},
		gitResolver:   gitResolver,
	}, nil
}
