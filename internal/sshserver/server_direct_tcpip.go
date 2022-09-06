package sshserver

import (
	"github.com/gliderlabs/ssh"
	"github.com/yolo-sh/agent/constants"
	"github.com/yolo-sh/agent/internal/docker"
	gossh "golang.org/x/crypto/ssh"
)

// proxyDirectTCPIPChannel is used to ovewrite dest addr
// with the container one during "direct-tcpip" channel
type proxyDirectTCPIPChannel struct {
	originalChannel gossh.NewChannel
}

func (ch *proxyDirectTCPIPChannel) Accept() (gossh.Channel, <-chan *gossh.Request, error) {
	return ch.originalChannel.Accept()
}

func (ch *proxyDirectTCPIPChannel) Reject(reason gossh.RejectionReason, message string) error {
	return ch.originalChannel.Reject(reason, message)
}

func (ch *proxyDirectTCPIPChannel) ChannelType() string {
	return ch.originalChannel.ChannelType()
}

func (ch *proxyDirectTCPIPChannel) ExtraData() []byte {
	msg := directTCPIPMsg{}
	err := gossh.Unmarshal(ch.originalChannel.ExtraData(), &msg)

	if err != nil {
		return ch.originalChannel.ExtraData()
	}

	dockerClient, err := docker.NewDefaultClient()

	if err != nil {
		return ch.originalChannel.ExtraData()
	}

	containerIPAddress, err := docker.LookupContainerIP(
		dockerClient,
		constants.DockerContainerName,
	)

	if err != nil || containerIPAddress == nil {
		return ch.originalChannel.ExtraData()
	}

	msg.DestAddr = *containerIPAddress
	return gossh.Marshal(msg)
}

// directTCPIPMsg is a struct used for SSH_MSG_CHANNEL_OPEN message
// with "direct-tcpip" string.
type directTCPIPMsg struct {
	DestAddr string
	DestPort uint32

	OriginAddr string
	OriginPort uint32
}

// handleDirectTCPIP is used to forward local conn to a remote port.
// Corresponds to the "direct-tcpip" channel type.
// Used by the local VSCode instance to reach code-server.
func handleDirectTCPIP(
	srv *ssh.Server,
	conn *gossh.ServerConn,
	newChan gossh.NewChannel,
	ctx ssh.Context,
) {
	proxyChannel := &proxyDirectTCPIPChannel{
		originalChannel: newChan,
	}

	ssh.DirectTCPIPHandler(
		srv,
		conn,
		proxyChannel,
		ctx,
	)
}
