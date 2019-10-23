package main

import (
	"time"

	tp "github.com/lazyweb/teleport"
)

//go:generate go build $GOFILE

func main() {
	defer tp.SetLoggerLevel("ERROR")()

	cli := tp.NewPeer(tp.PeerConfig{})
	defer cli.Close()
	cli.SetTLSConfig(tp.GenerateTLSConfigForClient())

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		tp.Fatalf("%v", stat)
	}

	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		tp.WithAddMeta("author", "henrylee2cn"),
	).Status()
	if !stat.OK() {
		tp.Fatalf("%v", stat)
	}
	tp.Printf("result: %d", result)
	tp.Printf("Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

// Push push handler
type Push struct {
	tp.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *tp.Status {
	tp.Printf("%s", *arg)
	return nil
}
