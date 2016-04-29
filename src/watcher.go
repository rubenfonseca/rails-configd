package src

import (
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Watcher is responsible for watching a directory structure on etcd
type Watcher struct {
	EtcdWatcher client.Watcher
	StopChannel chan bool
}

// Loop starts the watcher, returning new responses on the returned channel
func (w *Watcher) Loop() chan *client.Response {
	resultChannel := make(chan *client.Response)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			response, err := w.EtcdWatcher.Next(ctx)
			if err != nil {
				if err == context.Canceled {
					close(resultChannel)
					return
				}
				panic(err)
			}

			resultChannel <- response
		}
	}()

	go func() {
		<-w.StopChannel
		cancel()
	}()

	return resultChannel
}
