// Rails-configd - an online rails configuration generator using etcd data
//
// Standard usage:
//   (inside your Rails app)
//   $ rails-configd -etcd http://localhost:4001 -etcd-dir /rails/production -env production -renderer yaml -reloader touch
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/rubenfonseca/rails-configd/src"
)

var usageMessage = `rails-configd [%s]
This is a tool for watching over etcd tree, create config files for Rails, and restart the Rails processes.

Usage: %s [options]

The following options are recognized:
`

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage, releaseVersion, os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func loop(receiverChannel chan *etcd.Response, env src.Env) {
	for response := range receiverChannel {
		key := env.NakedKey(response.Node.Key, *env.EtcdDir)
		parts := strings.Split(key, "/")
		env.UpdateData(parts, response.Node.Value, response.Action, env.Data)

		log.Printf("[CHANGE]: %s %s %s", response.Action, key, response.Node.Value)

		env.Cycle()
	}
}

func main() {
	env := src.Env{}
	env.Data = make(map[string]interface{})

	env.Etcd = flag.String("etcd", "http://localhost:4001", "etcd address location")
	env.EtcdDir = flag.String("etcd-dir", "/rails_app01", "etcd directory that contains the configurations")

	rendererPtr := flag.String("renderer", "yaml", "The renderer to use when outputing the configs")
	reloaderPtr := flag.String("reloader", "touch", "The strategy to reload the Rails app")

	src.RegisterRendererFlags()
	src.RegisterReloaderFlags()

	flag.Usage = usage
	flag.Parse()

	// renderer
	renderer, err := src.OpenRenderer(*rendererPtr)
	if err != nil {
		panic(err)
	}
	env.Renderer = renderer

	// reloader
	env.Reloader, err = src.OpenReloader(*reloaderPtr)
	if err != nil {
		panic(err)
	}

	// etcd
	receiverChannel := make(chan *etcd.Response)
	stopChannel := make(chan bool)
	etcdClient := etcd.NewClient([]string{*env.Etcd})
	success := etcdClient.SyncCluster()
	if !success {
		log.Fatal("Cannot sync with etcd machines, please check --etcd")
	}

	etcdResponse, err := etcdClient.Get(*env.EtcdDir, false, true)
	if err != nil {
		panic(err)
	}
	if !etcdResponse.Node.Dir {
		panic("etc-dir should be a directory")
	}
	env.BuildData(*etcdResponse.Node, *env.EtcdDir, env.Data)
	env.Cycle()

	log.Printf("[MAIN] Waiting for changes from etcd @ %s", *env.EtcdDir)
	go etcdClient.Watch(*env.EtcdDir, 0, true, receiverChannel, stopChannel)

	// signals
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)
	go func() {
		for _ = range osSignal {
			log.Print("Interrupt received, finishing")
			stopChannel <- true
		}
	}()

	loop(receiverChannel, env)
}
