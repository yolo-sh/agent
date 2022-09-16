package env

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/yolo-sh/agent-container/constants"
	"github.com/yolo-sh/agent/internal/docker"
	"github.com/yolo-sh/agent/proto"
)

func EnsureDockerContainerRunning(
	dockerClient *client.Client,
	stream proto.Agent_BuildAndStartEnvServer,
	containerHostname string,
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
			"Pulling docker image (%s)",
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

	err = stream.Send(&proto.BuildAndStartEnvReply{
		WaitingForContainerAgent: true,
	})

	if err != nil {
		return err
	}

	createdDockerContainer, err := dockerClient.ContainerCreate(
		context.TODO(),

		&container.Config{
			WorkingDir: constants.WorkspaceDirPath,
			Image:      constants.DockerImageName,
			Entrypoint: constants.DockerContainerEntrypoint,
			Cmd:        strslice.StrSlice{},
			Hostname:   containerHostname,
			User:       "root",
		},

		&container.HostConfig{
			AutoRemove: false,
			Binds:      buildHostMounts(),
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
			Runtime: "sysbox-runc",
		},

		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"bridge": {
					IPAddress: constants.DockerContainerIPAddress,
				},
			},
		},

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

func WaitForAgentContainer() (
	returnedError error,
) {
	pollTimeoutChan := time.After(1 * time.Minute)
	pollSleepDuration := time.Second * 4
	connTimeout := time.Second * 4

	for {
		select {
		case <-pollTimeoutChan:
			return
		default:
			conn, err := net.DialTimeout(
				constants.GRPCServerAddrProtocol,
				constants.GRPCServerAddr,
				connTimeout,
			)

			// Make sure timeout returns last error
			returnedError = err

			if err != nil {
				break // wait pollSleepDuration and retry until timeout
			}

			conn.Close()
			return
		}

		time.Sleep(pollSleepDuration)
	}
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
		// Yolo configuration
		fmt.Sprintf(
			"%s:%s",
			constants.YoloConfigDirPath,
			constants.YoloConfigDirPath,
		),
	}
}
