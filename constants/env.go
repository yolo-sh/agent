package constants

const (
	YoloUserName                      = "yolo"
	YoloUserAuthorizedSSHKeysFilePath = "/home/yolo/.ssh/authorized_keys"

	DockerGroupName                   = "docker"
	DockerImageName                   = "yolosh/base-env:latest"
	DockerContainerName               = "yolo-env-container"
	DockerContainerEntrypointFilePath = "/yolo_entrypoint.sh"

	WorkspaceDirPath = "/home/yolo/workspace"

	WorkspaceConfigDirPath        = "/home/recode/.workspace-config"
	WorkspaceConfigFilePath       = WorkspaceConfigDirPath + "/recode.workspace"
	VSCodeWorkspaceConfigFilePath = WorkspaceConfigDirPath + "/recode.code-workspace"

	GitHubPublicSSHKeyFilePath = "/home/yolo/.ssh/yolo_github.pub"
	GitHubPublicGPGKeyFilePath = "/home/yolo/.gnupg/yolo_github_gpg_public.pgp"
)

var (
	DockerContainerStartCmd = []string{
		"sleep",
		"infinity",
	}
)
