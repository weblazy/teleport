package main

import (
	"time"

	tp "github.com/lazyweb/teleport"
	"github.com/lazyweb/teleport/plugin/proxy"
)

//go:generate go build $GOFILE

func main() {
	defer tp.FlushLogger()
	srv := tp.NewPeer(
		tp.PeerConfig{
			ListenPort: 8080,
		},
		newProxyPlugin(),
	)
	srv.ListenAndServe()
}

func newProxyPlugin() tp.Plugin {
	cli := tp.NewPeer(tp.PeerConfig{RedialTimes: 3})
	var sess tp.Session
	var stat *tp.Status
DIAL:
	sess, stat = cli.Dial(":9090")
	if !stat.OK() {
		tp.Warnf("%v", stat)
		time.Sleep(time.Second * 3)
		goto DIAL
	}
	return proxy.NewPlugin(func(*proxy.Label) proxy.Forwarder {
		return sess
	})
}
