package constants

const (
	GRPCServerAddrProtocol = "unix"
	GRPCServerAddr         = "/tmp/yolo_grpc.sock"

	SSHServerListenPort      = "2200"
	SSHServerListenAddr      = ":" + SSHServerListenPort
	SSHServerHostKeyFilePath = "/home/yolo/.ssh/yolo_ssh_server_host_key"

	InitInstanceScriptRepoPath = "yolo-sh/agent/internal/grpcserver/init_instance.sh"
)
