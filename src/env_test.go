package src

import (
	"testing"

	"github.com/bmizerany/assert"
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

func TestNakedKey(t *testing.T) {
	env := Env{}

	key := env.NakedKey("/rails/production/foo", "/rails/production")
	assert.Equal(t, key, "foo")

	key = env.NakedKey("/rails/production/foo/bar", "/rails/production")
	assert.Equal(t, key, "foo/bar")
}

func TestCycle(t *testing.T) {
	env := Env{Renderer: new(MockRenderer), Reloader: new(MockReloader)}

	env.Cycle()
	assert.Equal(t, env.Renderer.(*MockRenderer).Called, true)
	assert.Equal(t, env.Reloader.(*MockReloader).Called, true)
}
