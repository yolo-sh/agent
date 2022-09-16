package constants

import (
	agentContainerConsts "github.com/yolo-sh/agent-container/constants"
)

const (
	GRPCServerAddrProtocol = "unix"
	GRPCServerAddr         = agentContainerConsts.YoloConfigDirPath + "/agent-grpc.sock"

	SSHServerListenPort      = "2200"
	SSHServerListenAddr      = ":" + SSHServerListenPort
	SSHServerHostKeyFilePath = YoloUserHomeDirPath + "/.ssh/yolo-ssh-server-host-key"

	InitInstanceScriptRepoPath = "yolo-sh/agent/internal/grpcserver/init_instance.sh"
)
