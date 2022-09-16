package grpcserver

import (
	_ "embed"
	"fmt"
	"net"
	"os"

	"github.com/yolo-sh/agent/proto"
	"google.golang.org/grpc"
)

type agentServer struct {
	proto.UnimplementedAgentServer
}

func ListenAndServe(serverAddrProtocol, serverAddr string) error {
	tcpServer, err := net.Listen(serverAddrProtocol, serverAddr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if serverAddrProtocol == "unix" { // Make sure that the socket could be reached by the container agent
		if err := os.Chmod(serverAddr, 0660); err != nil {
			return fmt.Errorf("failed to set socket permissions: %v", err)
		}
	}

	grpcServer := grpc.NewServer()

	proto.RegisterAgentServer(grpcServer, &agentServer{})

	return grpcServer.Serve(tcpServer)
}
