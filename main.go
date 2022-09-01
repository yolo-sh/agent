package main

import (
	"log"
	"os"

	"github.com/yolo-sh/agent/constants"
	"github.com/yolo-sh/agent/internal/grpcserver"
	"github.com/yolo-sh/agent/internal/sshserver"
	"github.com/yolo-sh/agent/internal/system"
)

var (
	SSHServerAuthorizedUsers = []sshserver.AuthorizedUser{
		{
			UserName:               constants.YoloUserName,
			AuthorizedKeysFilePath: constants.YoloUserAuthorizedSSHKeysFilePath,
		},
	}
)

func main() {
	// Prevent "bind: address already in use" error
	err := ensureOldGRPCServerSocketRemoved(constants.GRPCServerAddr)

	if err != nil {
		log.Fatalf("%v", err)
	}

	go func() {
		log.Printf(
			"GRPC server listening at %s",
			constants.GRPCServerAddr,
		)

		err = grpcserver.ListenAndServe(
			constants.GRPCServerAddrProtocol,
			constants.GRPCServerAddr,
		)

		if err != nil {
			log.Fatalf("%v", err)
		}
	}()

	sshServerAuth := sshserver.NewAuth(
		system.NewFileManager(),
		sshserver.NewPrivateKeyManager(),
		sshserver.NewAuthorizedKeyManager(),
		constants.SSHServerHostKeyFilePath,
		SSHServerAuthorizedUsers,
	)

	sshServerBuilder := sshserver.NewServerBuilder(
		sshServerAuth,
		constants.SSHServerListenAddr,
	)

	sshServer, err := sshServerBuilder.Build()

	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf(
		"SSH server listening at %s",
		sshServer.Addr,
	)

	if err = sshServer.ListenAndServe(); err != nil {
		log.Fatalf("%v", err)
	}
}

func ensureOldGRPCServerSocketRemoved(socketPath string) error {
	return os.RemoveAll(socketPath)
}
