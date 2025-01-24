package main

import (
	"fmt"
	"net/netip"
	"time"

	netreactors "github.com/lance-e/net-reactors"
)

func connectionCallback(conn *netreactors.TcpConnection) {
	// conn.Send([]byte("已上线\n"))
}

func messageCallback(conn *netreactors.TcpConnection, buffer *netreactors.Buffer, t time.Time) {
	fmt.Printf("get message from server:%s\n", buffer.RetrieveAllString())
}

func HandleInput(conn *netreactors.TcpConnection) {
	for {

	}
}

func main() {
	ev := netreactors.NewEventLoop()
	addrport, _ := netip.ParseAddrPort("127.0.0.1:80")
	client := netreactors.NewTcpClient(ev, &addrport, "test")
	client.SetConnectionCallback(connectionCallback)
	client.SetMessageCallback(messageCallback)
	client.Connect()
	ev.Loop()
}
