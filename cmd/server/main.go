package main

import (
	"vchat/server"

	netreactors "github.com/lance-e/net-reactors"
)

func main() {
	// netreactors.Dlog.TurnOnLog()
	ev := netreactors.NewEventLoop()
	server.NewServer(ev, "127.0.0.1", 80, "chat", 0).Start()
	ev.Loop()
}
