package resolvers

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/custodian"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/modules"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	goModule "golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

type gitResolver struct {
	auth           transport.AuthMethod
	authMode       GitAuthMode
	moduleCacheDir string
}

func (f *gitResolver) Resolve(ctx context.Context, moduleIdentifier string) (custodian.Module, error) {
	resolvedIdentifier, err := f.getModule(moduleIdentifier)
	if err != nil {
		return nil, err
	}

	moduleDir := f.modulePathFromIdentifier(resolvedIdentifier)

	rFs := os.DirFS(moduleDir)
	return modules.NewModuleFromFS(resolvedIdentifier, rFs)
}

func (f *gitResolver) modulePathFromIdentifier(moduleIdentifier string) string {
	return path.Join(f.moduleCacheDir, moduleIdentifier)
}

func (f *gitResolver) getModule(moduleIdentifier string) (string, error) {
	// module already exists in module cache
	if utils.DirExists(f.modulePathFromIdentifier(moduleIdentifier)) {
		return moduleIdentifier, nil
	}

	// parse module identifier
	mId := GitModuleIdentifier(moduleIdentifier)
	remoteIdentifier, branch, version := mId.Remote(), mId.Branch(), mId.Version()
	// handle branch if present
	referenceName := plumbing.ReferenceName(branch)
	if branch != "" {
		referenceName = plumbing.NewBranchReferenceName(branch)
	}
	// clone the module
	log.Printf("Cloning: %s", moduleIdentifier)
	cloneOptions := &git.CloneOptions{
		URL:           buildRemoteURL(f.authMode, remoteIdentifier),
		Progress:      os.Stderr,
		Auth:          f.auth,
		ReferenceName: referenceName,
	}

	tmpDir, err := os.MkdirTemp("", "jnetx-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	repo, err := git.PlainClone(tmpDir, false, cloneOptions)
	if err != nil && err != git.ErrRepositoryAlreadyExists && !errors.Is(err, fs.ErrExist) {
		return "", err
	}

	// Find the complete version and commit hash
	completeVersion, commitHash, err := f.findPseudoVersion(version, repo)
	if err != nil {
		return "", err
	}

	// Get the worktree and checkout the specific commit
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(commitHash),
		Force: true,
	})

	// copy to module files removing the .git directory
	os.RemoveAll(path.Join(tmpDir, ".git"))
	moduleIdentifier = fmt.Sprintf("%s@%s", remoteIdentifier, completeVersion)
	targetDir := f.modulePathFromIdentifier(moduleIdentifier)
	os.RemoveAll(targetDir) // ensure target dir is clean
	err = os.CopyFS(targetDir, os.DirFS(tmpDir))
	log.Println("Module cloned", moduleIdentifier)

	return moduleIdentifier, err
}

func (f *gitResolver) findPseudoVersion(commitIdentifier string, repo *git.Repository) (string, string, error) {
	if commitIdentifier == "" {
		// If commitHashString is empty, get the latest commit hash from HEAD
		head, err := repo.Head()
		if err != nil {
			return "", "", err
		}
		commitIdentifier = head.Hash().String()
	} else if tagRef, err := repo.Tag(commitIdentifier); err == nil {
		// If commitHashString is a tag, resolve it to a commit hash
		tag, err := repo.Object(plumbing.AnyObject, tagRef.Hash())
		if err != nil {
			return "", "", err
		}

		commitIdentifier, err = resolveObjectToCommitHash(tag, repo)
		if err != nil {
			return "", "", err
		}
	} else if goModule.IsPseudoVersion(commitIdentifier) {
		// Handle pseudo-versions getting only the commit hash

		commitIdentifier, err = goModule.PseudoVersionRev(commitIdentifier)
		if err != nil {
			return "", "", err
		}
	}

	tagCommitMap, err := getCommitHashTagMap(repo)
	if err != nil {
		return "", "", err
	}

	gitLog, err := repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return "", "", err
	}
	defer gitLog.Close()

	var baseCommit *object.Commit

	for {
		commit, err := gitLog.Next()
		if err != nil {
			break
		}
		// look for the commit that matches the prefix
		if baseCommit == nil && strings.HasPrefix(commit.Hash.String(), commitIdentifier) {
			baseCommit = commit
			// Check if this commit has any semver tags
			if tagnames, ok := tagCommitMap[commit.Hash.String()]; ok {
				// Exact match with a semver tag
				semver.Sort(tagnames)
				return tagnames[0], commit.Hash.String(), nil
			}
		}
		// If we have found the base commit, look for the nearest semver tag in the history
		if baseCommit != nil {
			if tagnames, ok := tagCommitMap[commit.Hash.String()]; ok {
				semver.Sort(tagnames)
				pseudoVersion := goModule.PseudoVersion("", tagnames[0], baseCommit.Committer.When, commitIdentifier[0:12])
				return pseudoVersion, baseCommit.Hash.String(), nil
			}
		}
	}
	if baseCommit == nil {
		return "", "", fmt.Errorf("commit hash %s not found", commitIdentifier)
	}
	pseudoVersion := goModule.PseudoVersion("v0", "", baseCommit.Committer.When, commitIdentifier[0:12])
	return pseudoVersion, baseCommit.Hash.String(), nil
}

func NewGitResolver(targetDir string) (custodian.Resolver, error) {
	authMode, auth, err := getAuthMethodFromEnv()
	if err != nil {
		return nil, err
	}

	return &gitResolver{
		auth:           auth,
		authMode:       authMode,
		moduleCacheDir: targetDir,
	}, nil
}
