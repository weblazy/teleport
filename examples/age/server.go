package main

import (
	"context"
	"time"

	tp "github.com/weblazy/teleport"
)

//go:generate go build $GOFILE

func main() {
	defer tp.FlushLogger()
	srv := tp.NewPeer(tp.PeerConfig{
		PrintDetail:       true,
		CountTime:         true,
		ListenPort:        9090,
		DefaultSessionAge: time.Second * 7,
		DefaultContextAge: time.Second * 2,
	})
	srv.RouteCall(new(test))
	srv.ListenAndServe()
}

type test struct {
	tp.CallCtx
}

func (t *test) Ok(arg *string) (string, *tp.Status) {
	return *arg + " -> OK", nil
}

func (t *test) Timeout(arg *string) (string, *tp.Status) {
	tCtx, _ := context.WithTimeout(t.Context(), time.Second)
	time.Sleep(time.Second)
	select {
	case <-tCtx.Done():
		return "", tp.NewStatus(
			tp.CodeHandleTimeout,
			tp.CodeText(tp.CodeHandleTimeout),
			tCtx.Err().Error(),
		)
	default:
		return *arg + " -> Not Timeout", nil
	}
}

func (t *test) Break(*struct{}) (*struct{}, *tp.Status) {
	time.Sleep(time.Second * 3)
	select {
	case <-t.Session().CloseNotify():
		tp.Infof("the connection has gone away!")
		return nil, tp.NewStatus(tp.CodeConnClosed, "", "")
	default:
		return nil, nil
	}
}
