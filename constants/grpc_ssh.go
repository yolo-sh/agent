package constants

const (
	GRPCServerAddrProtocol = "unix"
	GRPCServerAddr         = "/tmp/yolo-grpc.sock"

	SSHServerListenPort      = "2200"
	SSHServerListenAddr      = ":" + SSHServerListenPort
	SSHServerHostKeyFilePath = "/home/" + YoloUserName + "/.ssh/yolo-ssh-server-host-key"

	InitInstanceScriptRepoPath = "yolo-sh/agent/internal/grpcserver/init_instance.sh"
)
