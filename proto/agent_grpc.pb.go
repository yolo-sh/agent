// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	InitInstance(ctx context.Context, in *InitInstanceRequest, opts ...grpc.CallOption) (Agent_InitInstanceClient, error)
	BuildAndStartEnv(ctx context.Context, in *BuildAndStartEnvRequest, opts ...grpc.CallOption) (Agent_BuildAndStartEnvClient, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) InitInstance(ctx context.Context, in *InitInstanceRequest, opts ...grpc.CallOption) (Agent_InitInstanceClient, error) {
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[0], "/agent.Agent/InitInstance", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentInitInstanceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Agent_InitInstanceClient interface {
	Recv() (*InitInstanceReply, error)
	grpc.ClientStream
}

type agentInitInstanceClient struct {
	grpc.ClientStream
}

func (x *agentInitInstanceClient) Recv() (*InitInstanceReply, error) {
	m := new(InitInstanceReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) BuildAndStartEnv(ctx context.Context, in *BuildAndStartEnvRequest, opts ...grpc.CallOption) (Agent_BuildAndStartEnvClient, error) {
	stream, err := c.cc.NewStream(ctx, &Agent_ServiceDesc.Streams[1], "/agent.Agent/BuildAndStartEnv", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentBuildAndStartEnvClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Agent_BuildAndStartEnvClient interface {
	Recv() (*BuildAndStartEnvReply, error)
	grpc.ClientStream
}

type agentBuildAndStartEnvClient struct {
	grpc.ClientStream
}

func (x *agentBuildAndStartEnvClient) Recv() (*BuildAndStartEnvReply, error) {
	m := new(BuildAndStartEnvReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentServer is the server API for Agent service.
// All implementations must embed UnimplementedAgentServer
// for forward compatibility
type AgentServer interface {
	InitInstance(*InitInstanceRequest, Agent_InitInstanceServer) error
	BuildAndStartEnv(*BuildAndStartEnvRequest, Agent_BuildAndStartEnvServer) error
	mustEmbedUnimplementedAgentServer()
}

// UnimplementedAgentServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (UnimplementedAgentServer) InitInstance(*InitInstanceRequest, Agent_InitInstanceServer) error {
	return status.Errorf(codes.Unimplemented, "method InitInstance not implemented")
}
func (UnimplementedAgentServer) BuildAndStartEnv(*BuildAndStartEnvRequest, Agent_BuildAndStartEnvServer) error {
	return status.Errorf(codes.Unimplemented, "method BuildAndStartEnv not implemented")
}
func (UnimplementedAgentServer) mustEmbedUnimplementedAgentServer() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_InitInstance_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InitInstanceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServer).InitInstance(m, &agentInitInstanceServer{stream})
}

type Agent_InitInstanceServer interface {
	Send(*InitInstanceReply) error
	grpc.ServerStream
}

type agentInitInstanceServer struct {
	grpc.ServerStream
}

func (x *agentInitInstanceServer) Send(m *InitInstanceReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Agent_BuildAndStartEnv_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BuildAndStartEnvRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AgentServer).BuildAndStartEnv(m, &agentBuildAndStartEnvServer{stream})
}

type Agent_BuildAndStartEnvServer interface {
	Send(*BuildAndStartEnvReply) error
	grpc.ServerStream
}

type agentBuildAndStartEnvServer struct {
	grpc.ServerStream
}

func (x *agentBuildAndStartEnvServer) Send(m *BuildAndStartEnvReply) error {
	return x.ServerStream.SendMsg(m)
}

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "agent.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InitInstance",
			Handler:       _Agent_InitInstance_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BuildAndStartEnv",
			Handler:       _Agent_BuildAndStartEnv_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "agent.proto",
}
