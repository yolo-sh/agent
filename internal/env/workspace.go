package env

import (
	"github.com/yolo-sh/agent/constants"
	"github.com/yolo-sh/agent/internal/system"
)

func PrepareWorkspace(repoOwner, repoName string) error {
	// The method "PrepareWorkspace" could
	// be called multiple times in case of error
	// so we need to make sure that our code is idempotent
	err := system.NewFileManager().RemoveDirContent(
		constants.WorkspaceDirPath,
	)

	if err != nil {
		return err
	}

	return cloneGitHubRepo(
		repoOwner,
		repoName,
		constants.WorkspaceDirPath,
	)
}
