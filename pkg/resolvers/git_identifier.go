package resolvers

import "strings"

// ModuleIdentifier represents a module identifier string
// e.g. host_fqdn/owner/repo[/branch]@version
//     |-----Git Remote-----|

const (
	VersionSeparator = "@"
)

type GitModuleIdentifier string

func (m GitModuleIdentifier) Remote() string {
	mIdData := strings.SplitN(string(m), VersionSeparator, 2)
	// remove branch if present
	remoteData := strings.SplitN(mIdData[0], "/", 4)
	if len(remoteData) == 4 {
		return strings.Join(remoteData[0:3], "/")
	}
	return mIdData[0]
}

func (m GitModuleIdentifier) Repo() string {
	mIdData := strings.SplitN(string(m), VersionSeparator, 2)
	remoteData := strings.SplitN(mIdData[0], "/", 4)
	if len(remoteData) >= 3 {
		return remoteData[2]
	}
	return ""
}

func (m GitModuleIdentifier) Version() string {
	mIdData := strings.SplitN(string(m), VersionSeparator, 2)
	if len(mIdData) == 2 {
		return mIdData[1]
	}
	return ""
}

func (m GitModuleIdentifier) Branch() string {
	remote := m.Remote()
	branchData := strings.SplitN(remote, "/", 4)
	if len(branchData) == 4 {
		return branchData[3]
	}
	return ""
}
