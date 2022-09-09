# Agent

This repository contains the source code of the Yolo Agent. 

The Yolo Agent is installed in the instance running your environment during its creation (most of the time via `cloud-init` but it may vary depending on the cloud provider used).

The main role of the Yolo Agent is to enable communication between your environments, the [CLI](https://github.com/yolo-sh/cli) and your code editor.

It is composed of two components: 

 - An `SSH server`.

 - A `gRPC server`.

## Table of contents
- [Requirements](#requirements)
- [Usage](#usage)
  - [Generating the gRPC server's code](#generating-the-grpc-servers-code)
- [Agent](#agent)
  - [SSH Server](#ssh-server)
  - [gRPC Server](#grpc-server)
- [License](#license)

## Requirements

The Yolo Agent only works on nix-based OS and requires:

  - `go >= 1.17`

  - `protoc >= 3.0` (see [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/))
  
  - `google.golang.org/protobuf/cmd/protoc-gen-go@latest` (install via `go install`)
  
  - `google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` (install via `go install`)

## Usage

The Yolo agent could be run using the `go run main.go` command. 

The `gRPC server` will listen on an Unix socket at `/tmp/yolo_grpc.sock` whereas the `SSH server` will listen on `:2200` by default.

### Generating the gRPC server's code

The `gRPC server`'s code could be generated by running the following command **in the `proto` directory**:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative agent.proto 
```

## Agent

The `SSH server` lets you access your environment via `SSH` without having to install and configure it. 

The `gRPC server`, on the other hand, is used to enable communication with the [Yolo CLI](https://github.com/yolo-sh/cli) (via `SSH`, using the OpenSSH's `Unix domain socket forwarding` feature).

<p align="center">
  <img src="https://user-images.githubusercontent.com/1233275/187863602-775b14db-f88d-4bfd-9b0b-c543643d020e.png" alt="infra" />
</p>

### SSH Server

The `SSH server` is the sole public-facing component. It will listen on port `2200` and redirect traffic to the `gRPC server` or to the `environment` depending on the requested `channel type`.

**The authentication will be done using the Public Key authentication method**. The key pair will be generated once, during the creation of the environment.

### gRPC server

The `gRPC server` listens on an Unix socket and, as a result, is not public-facing. It will be accessed by the [Yolo CLI](https://github.com/yolo-sh/cli) via `SSH`, using the OpenSSH's `Unix domain socket forwarding` feature.

It is principally used to build your environment as you can see in the service definition:

```proto
service Agent {
  rpc InitInstance (InitInstanceRequest) returns (stream InitInstanceReply) {}
  rpc BuildAndStartEnv (BuildAndStartEnvRequest) returns (stream BuildAndStartEnvReply) {}
}

message InitInstanceRequest {
  string env_name_slug = 1;
  string github_user_email = 2;
  string user_full_name = 3;
}

message InitInstanceReply {
  string log_line_header = 1;
  string log_line = 2;
  optional string github_ssh_public_key_content = 3;
  optional string github_gpg_public_key_content = 4;
}

message BuildAndStartEnvRequest {
  string env_repo_owner = 1;
  string env_repo_name = 2;
}

message BuildAndStartEnvReply {
  string log_line_header = 1;
  string log_line = 2;
}
```

The `InitInstance` method will run a [shell script](https://github.com/yolo-sh/agent/blob/main/internal/grpcserver/init_instance.sh) that will, among other things, install `Docker` and generate the `SSH` and `GPG` keys used in GitHub.

The `BuildAndStartEnv` method will clone your repositories and pull the `ghcr.io/yolo-sh/workspace-full` image.

**The two methods are idempotent**.

## License

Yolo is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
