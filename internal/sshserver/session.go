package sshserver

import (
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/yolo-sh/agent/constants"
	"github.com/yolo-sh/agent/internal/docker"
)

type SessionExecShellManager interface {
	ManageShellInEnv(sshSession ssh.Session) error
	ManageShellPTYInEnv(sshSession ssh.Session) error
	ManageExecInEnv(sshSession ssh.Session) error
	ManageShellPTY(sshSession ssh.Session) error
	ManageShell(sshSession ssh.Session) error
	ManageExec(sshSession ssh.Session) error
}

type Session struct {
	manager SessionExecShellManager
}

func NewSession(
	manager SessionExecShellManager,
) Session {

	return Session{
		manager: manager,
	}
}

func (s Session) Start(sshSession ssh.Session) {
	var sessionError error

	defer func() {
		if sessionError != nil {
			log.Println(sessionError)
			sshSession.Exit(1)
			return
		}

		sshSession.Exit(0)
	}()

	dockerClient, err := docker.NewDefaultClient()

	// We don't handle error here because
	// we want to be able to login to instance via SSH
	// even if the docker container cannot be reached
	if err != nil {
		log.Println(err)
	}

	isContainerRunning, err := docker.IsContainerRunning(
		dockerClient,
		constants.DockerContainerName,
	)

	// Same than previous comment
	if err != nil {
		log.Println(err)
	}

	if len(sshSession.Command()) == 0 { // "shell" session
		_, _, hasPTY := sshSession.Pty()

		if hasPTY {
			// if !isContainerRunning {
			// 	sessionError = s.manager.ManageShellPTY(sshSession)
			// 	return
			// }

			sessionError = s.manager.ManageShellPTYInEnv(sshSession)
			return
		}

		// if !isContainerRunning {
		// 	sessionError = s.manager.ManageShell(sshSession)
		// 	return
		// }

		sessionError = s.manager.ManageShellInEnv(sshSession)
		return
	}

	if !isContainerRunning {
		sessionError = s.manager.ManageExec(sshSession)
		return
	}

	// "exec" session
	sessionError = s.manager.ManageExecInEnv(sshSession)
}
