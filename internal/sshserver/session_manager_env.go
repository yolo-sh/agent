package sshserver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gliderlabs/ssh"
	"github.com/yolo-sh/agent-container/constants"
	"github.com/yolo-sh/agent/internal/docker"
	"github.com/yolo-sh/agent/internal/env"
)

func (s SessionManager) ManageShellInEnv(sshSession ssh.Session) error {
	_, _, hasPTY := sshSession.Pty()

	if hasPTY {
		return errors.New("expected no PTY, got PTY")
	}

	dockerClient, err := docker.NewDefaultClient()

	if err != nil {
		return err
	}

	workspaceConfig, err := env.LoadWorkspaceConfig(
		constants.WorkspaceConfigFilePath,
	)

	if err != nil {
		return err
	}

	exec, err := dockerClient.ContainerExecCreate(
		context.TODO(),
		constants.DockerContainerName,
		types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Detach:       false,
			Tty:          false,
			Cmd:          []string{"/bin/bash"},
			Env:          []string{},
			WorkingDir:   workspaceConfig.Repositories[0].RootDirPath,
			User:         constants.YoloUserName,
			Privileged:   true,
		},
	)

	if err != nil {
		return err
	}

	stream, err := dockerClient.ContainerExecAttach(
		context.TODO(),
		exec.ID,
		types.ExecStartCheck{},
	)

	if err != nil {
		return err
	}

	defer stream.Close()

	stdinChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(stream.Conn, sshSession)
		stdinChan <- err
	}()

	stdoutChan := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(
			sshSession,
			sshSession.Stderr(),
			stream.Reader,
		)

		stdoutChan <- err
	}()

	select {
	case stdoutErr := <-stdoutChan:
		return stdoutErr
	case stdinErr := <-stdinChan:
		return stdinErr
	}
}

func (s SessionManager) ManageShellPTYInEnv(sshSession ssh.Session) error {
	ptyReq, windowChan, hasPTY := sshSession.Pty()

	if !hasPTY {
		return errors.New("expected PTY, got no PTY")
	}

	dockerClient, err := docker.NewDefaultClient()

	if err != nil {
		return err
	}

	workspaceConfig, err := env.LoadWorkspaceConfig(
		constants.WorkspaceConfigFilePath,
	)

	if err != nil {
		return err
	}

	exec, err := dockerClient.ContainerExecCreate(
		context.TODO(),
		constants.DockerContainerName,
		types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Detach:       false,
			Tty:          true,
			Cmd: []string{
				"/bin/bash",
				"-c",
				fmt.Sprintf(
					// Display Ubuntu motd and run default shell for user
					"echo '' && for i in /etc/update-motd.d/*; do $i; done && echo '' && $(getent passwd %s | cut -d ':' -f 7)",
					constants.YoloUserName,
				),
			},
			Env: []string{
				fmt.Sprintf("TERM=%s", ptyReq.Term),
			},
			WorkingDir: workspaceConfig.Repositories[0].RootDirPath,
			User:       constants.YoloUserName,
			Privileged: true,
		},
	)

	if err != nil {
		return err
	}

	stream, err := dockerClient.ContainerExecAttach(
		context.TODO(),
		exec.ID,
		types.ExecStartCheck{
			Detach: false,
			Tty:    true,
		},
	)

	if err != nil {
		return err
	}

	defer stream.Close()

	resizeChan := make(chan error, 1)

	go func() {
		for window := range windowChan {
			err := dockerClient.ContainerExecResize(
				context.TODO(),
				exec.ID,
				types.ResizeOptions{
					Height: uint(window.Height),
					Width:  uint(window.Width),
				},
			)

			if err != nil {
				resizeChan <- err
				break
			}
		}
	}()

	stdinChan := make(chan error, 1)

	go func() {
		_, err := io.Copy(stream.Conn, sshSession)

		stdinChan <- err
	}()

	stdoutChan := make(chan error, 1)

	go func() {
		_, err := io.Copy(
			sshSession,
			stream.Reader,
		)

		stdoutChan <- err
	}()

	select {
	case resizeErr := <-resizeChan:
		return resizeErr
	case stdoutErr := <-stdoutChan:
		return stdoutErr
	case stdinErr := <-stdinChan:
		return stdinErr
	}
}

func (s SessionManager) ManageExecInEnv(sshSession ssh.Session) error {
	passedCmd := sshSession.Command()

	if len(passedCmd) == 0 {
		return errors.New("expected command, got nothing")
	}

	dockerClient, err := docker.NewDefaultClient()

	if err != nil {
		return err
	}

	workspaceConfig, err := env.LoadWorkspaceConfig(
		constants.WorkspaceConfigFilePath,
	)

	if err != nil {
		return err
	}

	exec, err := dockerClient.ContainerExecCreate(
		context.TODO(),
		constants.DockerContainerName,
		types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Detach:       false,
			Tty:          false,
			Cmd:          passedCmd,
			Env:          []string{},
			WorkingDir:   workspaceConfig.Repositories[0].RootDirPath,
			User:         constants.YoloUserName,
			Privileged:   true,
		},
	)

	if err != nil {
		return err
	}

	stream, err := dockerClient.ContainerExecAttach(
		context.TODO(),
		exec.ID,
		types.ExecStartCheck{},
	)

	if err != nil {
		return err
	}

	defer stream.Close()

	stdinChan := make(chan error, 1)

	go func() {
		_, err := io.Copy(stream.Conn, sshSession)

		stdinChan <- err
	}()

	stdoutChan := make(chan error, 1)

	go func() {
		_, err := stdcopy.StdCopy(
			sshSession,
			sshSession.Stderr(),
			stream.Reader,
		)

		stdoutChan <- err
	}()

	select {
	case stdoutErr := <-stdoutChan:
		if stdoutErr != nil {
			return stdoutErr
		}
	case stdinErr := <-stdinChan:
		if stdinErr != nil {
			return stdinErr
		}
	}

	containerInspect, err := dockerClient.ContainerExecInspect(
		context.TODO(),
		exec.ID,
	)

	if err != nil {
		return err
	}

	if containerInspect.ExitCode != 0 {
		return fmt.Errorf(
			"the command \"%s\" has returned a non-zero (%d) exit code",
			passedCmd,
			containerInspect.ExitCode,
		)
	}

	return nil
}
