package constants

const (
	YoloUserName                      = "yolo"
	YoloUserHomeDirPath               = "/home/" + YoloUserName
	YoloUserAuthorizedSSHKeysFilePath = YoloUserHomeDirPath + "/.ssh/authorized_keys"

	DockerImageTag  = "0.0.1"
	DockerImageName = "ghcr.io/yolo-sh/workspace-full:" + DockerImageTag

	DockerContainerName               = "yolo-env-container"
	DockerContainerEntrypointFilePath = "/entrypoint.sh"

	WorkspaceDirPath = YoloUserHomeDirPath + "/workspace"

	WorkspaceConfigDirPath        = YoloUserHomeDirPath + "/.workspace-config"
	WorkspaceConfigFilePath       = WorkspaceConfigDirPath + "/default.workspace"
	VSCodeWorkspaceConfigFilePath = WorkspaceConfigDirPath + "/default.code-workspace"

	GitHubPublicSSHKeyFilePath = YoloUserHomeDirPath + "/.ssh/" + YoloUserName + "_github.pub"
	GitHubPublicGPGKeyFilePath = YoloUserHomeDirPath + "/.gnupg/" + YoloUserName + "_github_gpg_public.pgp"
)

var (
	DockerContainerStartCmd = []string{
		"sleep",
		"infinity",
	}
)
