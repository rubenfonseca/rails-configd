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

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
	"github.com/rubenfonseca/rails-configd/src"
)

var (
	etcdFlag     = flag.String("etcd", "http://localhost:4001", "etcd address location")
	etcdDirFlag  = flag.String("etcd-dir", "/rails_app01", "etcd directory that contains the configurations")
	rendererFlag = flag.String("renderer", "yaml", "The renderer to use when outputing the configs")
	reloaderFlag = flag.String("reloader", "touch", "The strategy to reload the Rails app")
)

const usageMessage = `rails-configd [%s]
This is a tool for watching over etcd tree, create config files for Rails, and restart the Rails processes.

Usage: %s [options]

The following options are recognized:
`

func printHelp() {
	fmt.Fprintf(os.Stderr, usageMessage, releaseVersion, os.Args[0])
	flag.PrintDefaults()
}

func loop(receiverChannel chan *client.Response, env src.Env) {
	for response := range receiverChannel {
		key := env.NakedKey(response.Node.Key, *env.EtcdDir)
		parts := strings.Split(key, "/")
		env.UpdateData(parts, response.Node.Value, response.Action, env.Data)

		log.Printf("[CHANGE]: %s %s %s", response.Action, key, response.Node.Value)

		env.Cycle()
	}
}

func main() {
	var err error
	env := src.Env{Etcd: etcdFlag, EtcdDir: etcdDirFlag, Data: make(map[string]interface{})}

	src.RegisterRendererFlags()
	src.RegisterReloaderFlags()

	flag.Usage = printHelp
	flag.Parse()

	// renderer
	env.Renderer, err = src.OpenRenderer(*rendererFlag)
	if err != nil {
		panic(err)
	}

	// reloader
	env.Reloader, err = src.OpenReloader(*reloaderFlag)
	if err != nil {
		panic(err)
	}

	// main app context
	ctx, cancel := context.WithCancel(context.Background())

	// etcd
	etcdClientConfig := client.Config{Endpoints: []string{*env.Etcd}}
	etcdClient, err := client.New(etcdClientConfig)
	if err != nil {
		log.Fatal("Cannot connect to etcd machines, pleace check --etcd")
	}
	etcdKeyClient := client.NewKeysAPI(etcdClient)

	etcdGetOptions := &client.GetOptions{Recursive: true, Sort: false}
	etcdResponse, err := etcdKeyClient.Get(ctx, *env.EtcdDir, etcdGetOptions)
	if err != nil {
		panic(err)
	}
	if !etcdResponse.Node.Dir {
		panic("etc-dir should be a directory")
	}
	env.BuildData(*etcdResponse.Node, *env.EtcdDir, env.Data)
	env.Cycle()

	// watcher
	log.Printf("[MAIN] Waiting for changes from etcd @ %s", *env.EtcdDir)
	etcdWatcherOptions := &client.WatcherOptions{AfterIndex: 0, Recursive: true}
	etcdWatcher := etcdKeyClient.Watcher(*env.EtcdDir, etcdWatcherOptions)

	watcher := src.Watcher{EtcdWatcher: etcdWatcher}
	receiverChannel := watcher.Loop(ctx)

	// signals
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)
	go func() {
		for _ = range osSignal {
			log.Print("Interrupt received, finishing")
			cancel()
		}
	}()

	loop(receiverChannel, env)
}
