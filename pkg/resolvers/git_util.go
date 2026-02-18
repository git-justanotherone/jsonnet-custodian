package resolvers

import (
	"fmt"
	"os"
	"strings"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	ENV_PREFIX         = "CUSTODIAN_"
	ENV_GIT_AUTH_MODE  = ENV_PREFIX + "GIT_AUTH_MODE"
	ENV_GIT_SSH_KEY    = ENV_PREFIX + "GIT_SSH_KEY"
	ENV_GIT_AUTH_TOKEN = ENV_PREFIX + "GIT_AUTH_TOKEN"
	ENV_GIT_USER       = ENV_PREFIX + "GIT_USER"
	ENV_GIT_PASS       = ENV_PREFIX + "GIT_PASS"
	ENV_FILE_SUFFIX    = "_FILE"
)

type GitAuthMode string

const (
	GitAuthModeNone      GitAuthMode = "none"
	GitAuthModeAuthToken GitAuthMode = "auth-token"
	GitAuthModeBasicAuth GitAuthMode = "basic-auth"
	GitAuthModeSshKey    GitAuthMode = "ssh-key"
	GitAuthModeSshAgent  GitAuthMode = "ssh-agent"
)

func getAuthMethodFromEnv() (GitAuthMode, transport.AuthMethod, error) {
	switch GitAuthMode(os.Getenv(ENV_GIT_AUTH_MODE)) {
	case GitAuthModeAuthToken:
		return GitAuthModeAuthToken, &http.TokenAuth{
			Token: utils.GetEnvOrEmpty(ENV_GIT_AUTH_TOKEN),
		}, nil
	case GitAuthModeBasicAuth:
		return GitAuthModeBasicAuth, &http.BasicAuth{
			Username: utils.GetEnvOrEmpty(ENV_GIT_USER),
			Password: utils.GetEnvOrEmpty(ENV_GIT_PASS),
		}, nil
	case GitAuthModeSshKey:
		auth, err := ssh.NewPublicKeysFromFile(
			utils.GetEnvOrEmpty(ENV_GIT_USER), utils.GetEnvOrEmpty(ENV_GIT_SSH_KEY), utils.GetEnvOrEmpty(ENV_GIT_PASS),
		)
		return GitAuthModeSshKey, auth, err
	case GitAuthModeSshAgent:
		auth, err := ssh.NewSSHAgentAuth("git")
		return GitAuthModeSshAgent, auth, err
	default:
		return GitAuthModeNone, nil, nil
	}

}

func ParseModuleIdentifier(moduleIdentifier string) (string, string) {
	mIdData := strings.SplitN(moduleIdentifier, VersionSeparator, 2)
	if len(mIdData) == 2 {
		return mIdData[0], mIdData[1]
	}
	return mIdData[0], ""
}

func buildRemoteURL(authMethod GitAuthMode, remoteIdentifier string) string {
	switch authMethod {
	case GitAuthModeSshKey, GitAuthModeSshAgent:
		return "git@" + strings.Replace(remoteIdentifier, "/", ":", 1) // convert host_fqdn/owner/repo to git@host_fqdn:owner/repo
	default:
		return "https://" + remoteIdentifier // convert host_fqdn/owner/repo to https://host_fqdn/owner/repo
	}
}

func resolveObjectToCommitHash(gitObject object.Object, repo *git.Repository) (string, error) {
	switch obj := gitObject.(type) {
	case *object.Commit:
		return obj.Hash.String(), nil
	case *object.Tag:
		// Annotated tag, resolve to the target object
		target, err := repo.Object(plumbing.AnyObject, obj.Target)
		if err != nil {
			return "", err
		}
		return resolveObjectToCommitHash(target, repo)
	default:
		return "", fmt.Errorf("unsupported object type: %T", obj)
	}
}

func getCommitHashTagMap(repo *git.Repository) (map[string][]string, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	tagCommitMap := make(map[string][]string) // commit hash -> tag name
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		obj, err := repo.Object(plumbing.AnyObject, ref.Hash())
		if err != nil {
			return err
		}
		commitHash, err := resolveObjectToCommitHash(obj, repo)
		if err != nil {
			return err
		}
		tagCommitMap[commitHash] = append(tagCommitMap[commitHash], ref.Name().Short())
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tagCommitMap, nil
}
