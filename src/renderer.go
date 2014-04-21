package src

import (
	"fmt"
)

type Renderer interface {
	Render(env Env)
	RegisterFlags()
}

var renderers = make(map[string]Renderer)

func RegisterRenderer(name string, renderer Renderer) {
	if renderer == nil {
		panic("renderer: Register renderer is nil")
	}
	if _, dup := renderers[name]; dup {
		panic("renderer: Register called twice for renderer " + name)
	}
	renderers[name] = renderer
}

func OpenRenderer(rendererName string) (Renderer, error) {
	renderer, ok := renderers[rendererName]
	if !ok {
		return nil, fmt.Errorf("renderer: unkown driver %q (forgotten import?)", rendererName)
	}

	return renderer, nil
}

func RegisterRendererFlags() {
	for _, renderer := range renderers {
		renderer.RegisterFlags()
	}
}
