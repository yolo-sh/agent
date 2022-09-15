package sshserver

import (
	"github.com/gliderlabs/ssh"
	"github.com/yolo-sh/agent-container/constants"
	gossh "golang.org/x/crypto/ssh"
)

// proxyDirectTCPIPChannel is used during port forwarding
// to overwrite the destination IP with the env container one.
// "direct-tcpip" is for "client-to-server forwarded connections".
type proxyDirectTCPIPChannel struct {
	originalChannel gossh.NewChannel
}

func (ch *proxyDirectTCPIPChannel) Accept() (
	gossh.Channel,
	<-chan *gossh.Request,
	error,
) {

	return ch.originalChannel.Accept()
}

func (ch *proxyDirectTCPIPChannel) Reject(
	reason gossh.RejectionReason,
	message string,
) error {

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

	msg.DestAddr = constants.DockerContainerIPAddress
	return gossh.Marshal(msg)
}

// directTCPIPMsg represents the message
// sent during the opening of "direct-tcpip" channels
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
