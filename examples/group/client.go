package main

import (
	tp "github.com/lazyweb/teleport"
)

//go:generate go build $GOFILE

func main() {
	defer tp.SetLoggerLevel("ERROR")()

	cli := tp.NewPeer(
		tp.PeerConfig{},
	)
	defer cli.Close()
	group := cli.SubRoute("/cli")
	group.RoutePush(new(push))

	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		tp.Fatalf("%v", stat)
	}

	var result int
	stat = sess.Call("/srv/math/v2/add_2",
		[]int{1, 2, 3, 4, 5},
		&result,
		tp.WithSetMeta("push_status", "yes"),
	).Status()

	if !stat.OK() {
		tp.Fatalf("%v", stat)
	}
	tp.Printf("result: %d", result)
}

type push struct {
	tp.PushCtx
}

func (p *push) ServerStatus(arg *string) *tp.Status {
	tp.Printf("server status: %s", *arg)
	return nil
}
