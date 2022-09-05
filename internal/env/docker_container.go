package env

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/yolo-sh/agent/constants"
	"github.com/yolo-sh/agent/internal/docker"
	"github.com/yolo-sh/agent/proto"
)

func EnsureDockerContainerRunning(
	dockerClient *client.Client,
	stream proto.Agent_BuildAndStartEnvServer,
) error {

	isContainerRunning, err := docker.IsContainerRunning(
		dockerClient,
		constants.DockerContainerName,
	)

	if err != nil {
		return err
	}

	if isContainerRunning {
		return nil
	}

	dockerContainer, err := docker.LookupContainer(
		dockerClient,
		constants.DockerContainerName,
	)

	if err != nil {
		return err
	}

	if dockerContainer != nil { // Container exists but is not running
		return dockerClient.ContainerStart(
			context.TODO(),
			dockerContainer.ID,
			types.ContainerStartOptions{},
		)
	}

	err = stream.Send(&proto.BuildAndStartEnvReply{
		LogLineHeader: fmt.Sprintf(
			"Pulling Docker image %s",
			constants.DockerImageName,
		),
	})

	if err != nil {
		return err
	}

	err = pullBaseEnvImage(
		dockerClient,
		stream,
	)

	if err != nil {
		return err
	}

	createdDockerContainer, err := dockerClient.ContainerCreate(
		context.TODO(),

		&container.Config{
			WorkingDir: constants.WorkspaceDirPath,
			Image:      constants.DockerImageName,
			User:       constants.YoloUserName,
			Entrypoint: strslice.StrSlice{
				constants.DockerContainerEntrypointFilePath,
			},
			Cmd: constants.DockerContainerStartCmd,
		},

		&container.HostConfig{
			AutoRemove:  false,
			Binds:       buildHostMounts(),
			//NetworkMode: container.NetworkMode("host"),
			Privileged:  true,
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
			Runtime: "sysbox-runc",
		},

		nil,

		nil,

		constants.DockerContainerName,
	)

	if err != nil {
		return err
	}

	return dockerClient.ContainerStart(
		context.TODO(),
		createdDockerContainer.ID,
		types.ContainerStartOptions{},
	)
}

func EnsureDockerContainerRemoved(dockerClient *client.Client) error {
	dockerContainer, err := docker.LookupContainer(
		dockerClient,
		constants.DockerContainerName,
	)

	if err != nil {
		return err
	}

	if dockerContainer == nil {
		return nil
	}

	return dockerClient.ContainerRemove(
		context.TODO(),
		dockerContainer.ID,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)
}

func pullBaseEnvImage(
	dockerClient *client.Client,
	stream proto.Agent_BuildAndStartEnvServer,
) error {

	reader, err := dockerClient.ImagePull(
		context.TODO(),
		constants.DockerImageName,
		types.ImagePullOptions{},
	)

	if err != nil {
		return err
	}

	defer reader.Close()

	return docker.HandlePullOutput(
		reader,
		func(logLine string) error {
			return stream.Send(&proto.BuildAndStartEnvReply{
				LogLine: logLine,
			})
		},
	)
}

func buildHostMounts() []string {
	return []string{
		// Working dir

		fmt.Sprintf(
			"%s:%s",
			constants.WorkspaceDirPath,
			constants.WorkspaceDirPath,
		),

		fmt.Sprintf(
			"%s:%s",
			constants.WorkspaceConfigDirPath,
			constants.WorkspaceConfigDirPath,
		),

		/* Config files are mounted to /etc/
		   to let users overwrite them, if needed,
		   using config files in home dir. */

		// Git config

		fmt.Sprintf(
			"/home/%s/.gitconfig:/etc/gitconfig",
			constants.YoloUserName,
		),

		// SSH config

		fmt.Sprintf(
			"/home/%s/.ssh/config:/etc/ssh/ssh_config",
			constants.YoloUserName,
		),

		fmt.Sprintf(
			"/home/%s/.ssh/known_hosts:/etc/ssh/ssh_known_hosts",
			constants.YoloUserName,
		),

		// SSH GitHub keys

		fmt.Sprintf(
			"/home/%s/.ssh/yolo_github:/home/%s/.ssh/yolo_github",
			constants.YoloUserName,
			constants.YoloUserName,
		),

		fmt.Sprintf(
			"/home/%s/.ssh/yolo_github.pub:/home/%s/.ssh/yolo_github.pub",
			constants.YoloUserName,
			constants.YoloUserName,
		),

		// GnuPG GitHub keys

		fmt.Sprintf(
			"/home/%s/.gnupg/yolo_github_gpg_public.pgp:/home/%s/.gnupg/yolo_github_gpg_public.pgp",
			constants.YoloUserName,
			constants.YoloUserName,
		),

		fmt.Sprintf(
			"/home/%s/.gnupg/yolo_github_gpg_private.pgp:/home/%s/.gnupg/yolo_github_gpg_private.pgp",
			constants.YoloUserName,
			constants.YoloUserName,
		),

		// Docker daemon socket

		// "/var/run/docker.sock:/var/run/docker.sock",
	}
}
