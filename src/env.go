package src

import (
	"log"
	"strings"

	"github.com/coreos/etcd/client"
)

// Env represents all the necessary data the core needs to run
type Env struct {
	// Etcd address
	Etcd *string
	// Directory inside etcd that contains the configuration
	EtcdDir *string
	// Structure that holds the configuration data in memory
	Data map[string]interface{}
	// An instance of a renderer
	Renderer Renderer
	// An instance of a reloader
	Reloader Reloader
}

// Cycles the rails environemnt, by rendering a new configuration
// file and reloading the Rails processes. Uses the existing renderer
// and reloader instances.
func (env *Env) Cycle() {
	log.Printf("[ENV] Rendering and reloading...")

	env.Renderer.Render(*env)
	env.Reloader.Reload()
}

// Taking a etcd node and a prefix, updates the in memory data.
// If the etcd node represents a nested directory, this function calls recursively
// with the new prefix, trying to create a tree structure in memory.
func (env *Env) BuildData(node client.Node, prefix string, data map[string]interface{}) {
	for i := range node.Nodes {
		node := node.Nodes[i]
		key := env.NakedKey(node.Key, prefix)

		if node.Dir {
			data[key] = make(map[string]interface{})
			env.BuildData(*node, prefix+"/"+key, data[key].(map[string]interface{}))
		} else {
			data[key] = node.Value
		}
	}
}

// Updates the data from an etcd watch update. Takes into consideration the type of action
// (set or delete) and navigates through the parts until if finds the correct node to update.
func (env *Env) UpdateData(parts []string, value string, action string, data map[string]interface{}) {
	head := parts[0]
	tail := parts[1:]

	if len(tail) == 0 {
		if action == "set" {
			data[head] = value
		}
		if action == "delete" {
			delete(data, head)
		}
	} else {
		if _, ok := data[head]; !ok {
			newData := make(map[string]interface{})
			data[head] = newData
		}
		env.UpdateData(tail, value, action, data[head].(map[string]interface{}))
	}
}

// Removes the prefix from a key, including trailing slashes
func (env *Env) NakedKey(key string, prefix string) string {
	key = strings.Replace(key, prefix, "", -1)
	return strings.Replace(key, "/", "", 1)
}
