package grpcserver

import (
	"github.com/yolo-sh/agent/internal/docker"
	"github.com/yolo-sh/agent/internal/env"
	"github.com/yolo-sh/agent/proto"
)

func (s *agentServer) BuildAndStartEnv(
	req *proto.BuildAndStartEnvRequest,
	stream proto.Agent_BuildAndStartEnvServer,
) error {

	dockerClient, err := docker.NewDefaultClient()

	if err != nil {
		return err
	}

	// The method "BuildAndStartEnv" may be run multiple times
	// so we need to ensure idempotency
	err = env.EnsureDockerContainerRemoved(dockerClient)

	if err != nil {
		return err
	}

	err = env.PrepareWorkspace(
		req.EnvRepoOwner,
		req.EnvRepoName,
	)

	if err != nil {
		return err
	}

	return env.EnsureDockerContainerRunning(
		dockerClient,
		stream,
	)
}
