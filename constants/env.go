package constants

const (
	YoloUserName                      = "yolo"
	YoloUserAuthorizedSSHKeysFilePath = "/home/yolo/.ssh/authorized_keys"

	DockerGroupName                   = "docker"
	DockerImageName                   = "yolosh/base-env:0.0.2-dev"
	DockerContainerName               = "yolo-env-container"
	DockerContainerEntrypointFilePath = "/yolo_entrypoint.sh"

	WorkspaceDirPath = "/home/yolo/workspace"

	WorkspaceConfigDirPath        = "/home/yolo/.workspace-config"
	WorkspaceConfigFilePath       = WorkspaceConfigDirPath + "/yolo.workspace"
	VSCodeWorkspaceConfigFilePath = WorkspaceConfigDirPath + "/yolo.code-workspace"

	GitHubPublicSSHKeyFilePath = "/home/yolo/.ssh/yolo_github.pub"
	GitHubPublicGPGKeyFilePath = "/home/yolo/.gnupg/yolo_github_gpg_public.pgp"
)

var (
	DockerContainerStartCmd = []string{
		"/sbin/init",
		"--log-level=err",
	}
)
