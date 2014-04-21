package src

import (
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

type Env struct {
	RailsEnv *string
	Etcd     *string
	EtcdDir  *string
	Data     map[string]interface{}
	Renderer Renderer
	Reloader Reloader
}

func (env *Env) Cycle() {
	log.Printf("[ENV] Rendering and reloading...")

	env.Renderer.Render(*env)
	env.Reloader.Reload()
}

func (env *Env) BuildData(node etcd.Node, prefix string, data map[string]interface{}) {
	for i := range node.Nodes {
		node := node.Nodes[i]
		key := env.NakedKey(node.Key, prefix)

		if node.Dir {
			data[key] = make(map[string]interface{})
			env.BuildData(node, prefix+"/"+key, data[key].(map[string]interface{}))
		} else {
			data[key] = node.Value
		}
	}
}

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

func (env *Env) NakedKey(key string, prefix string) string {
	key = strings.Replace(key, prefix, "", -1)
	return strings.Replace(key, "/", "", 1)
}
