package grpcserver

import (
	"context"
	"io"

	"github.com/yolo-sh/agent-container/constants"
	agentContainerProto "github.com/yolo-sh/agent-container/proto"
	"github.com/yolo-sh/agent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (*agentServer) InitEnv(
	req *proto.InitEnvRequest,
	stream proto.Agent_InitEnvServer,
) error {

	grpcConn, err := grpc.Dial(
		constants.GRPCServerUri,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return err
	}

	defer grpcConn.Close()

	agentContainerClient := agentContainerProto.NewAgentClient(grpcConn)

	initStream, err := agentContainerClient.Init(
		context.TODO(),
		&agentContainerProto.InitRequest{
			EnvRepoOwner:         req.EnvRepoOwner,
			EnvRepoName:          req.EnvRepoName,
			EnvRepoLanguagesUsed: req.EnvRepoLanguagesUsed,
			UserFullName:         req.UserFullName,
			GithubUserEmail:      req.GithubUserEmail,
		},
	)

	if err != nil {
		return err
	}

	for {
		initReply, err := initStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		err = stream.Send(&proto.InitEnvReply{
			LogLineHeader:             initReply.LogLineHeader,
			LogLine:                   initReply.LogLine,
			GithubSshPublicKeyContent: initReply.GithubSshPublicKeyContent,
			GithubGpgPublicKeyContent: initReply.GithubGpgPublicKeyContent,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
