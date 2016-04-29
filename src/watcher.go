package src

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Watcher is responsible for watching a directory structure on etcd
type Watcher struct {
	EtcdWatcher client.Watcher
}

// Loop starts the watcher, returning new responses on the returned channel
func (w *Watcher) Loop(ctx context.Context) chan *client.Response {
	resultChannel := make(chan *client.Response)

	go watchLoop(w.EtcdWatcher, ctx, resultChannel)
	return resultChannel
}

func watchLoop(w client.Watcher, ctx context.Context, r chan *client.Response) {
	for {
		response, err := w.Next(ctx)
		if err == context.Canceled {
			close(r)
			return
		}
		if err != nil {
			panic(err)
		}

		r <- response
	}
}
