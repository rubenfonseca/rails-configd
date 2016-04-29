package src

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/coreos/etcd/client"
)

type MockRenderer struct {
	Called bool
}

func (r *MockRenderer) Render(env Env) {
	r.Called = true
}
func (r *MockRenderer) RegisterFlags() {
}

type MockReloader struct {
	Called bool
}

func (r *MockReloader) Reload() {
	r.Called = true
}
func (r *MockReloader) RegisterFlags() {
}

func TestCycle(t *testing.T) {
	env := Env{Renderer: new(MockRenderer), Reloader: new(MockReloader)}

	env.Cycle()
	assert.Equal(t, env.Renderer.(*MockRenderer).Called, true)
	assert.Equal(t, env.Reloader.(*MockReloader).Called, true)
}

func TestBuildData(t *testing.T) {
	env := Env{}

	hostnameNode := client.Node{Key: "/rails/mongodb/hostname", Value: "localhost"}
	mongoDbNode := client.Node{Key: "/rails/mongodb", Dir: true, Nodes: client.Nodes{&hostnameNode}}
	dirNode := client.Node{Dir: true, Nodes: client.Nodes{&mongoDbNode}}

	data := map[string]interface{}{}
	env.BuildData(dirNode, "/rails", data)

	mongodb := data["mongodb"].(map[string]interface{})
	assert.Equal(t, mongodb["hostname"], "localhost")
}

func TestUpdateData(t *testing.T) {
	env := Env{}

	data := map[string]interface{}{"mongodb": map[string]interface{}{"hostname": "localhost"}}

	env.UpdateData([]string{"mongodb", "hostname"}, "google.com", "set", data)

	mongodb := data["mongodb"].(map[string]interface{})
	assert.Equal(t, mongodb["hostname"], "google.com")

	env.UpdateData([]string{"mongodb", "hostname"}, "", "delete", data)
	mongodb = data["mongodb"].(map[string]interface{})
	assert.Equal(t, mongodb["hostname"], nil)
}

func TestNakedKey(t *testing.T) {
	env := Env{}

	key := env.NakedKey("/rails/production/foo", "/rails/production")
	assert.Equal(t, key, "foo")

	key = env.NakedKey("/rails/production/foo/bar", "/rails/production")
	assert.Equal(t, key, "foo/bar")
}
