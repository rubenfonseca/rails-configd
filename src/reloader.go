package src

import "fmt"

type Reloader interface {
	Reload()
	RegisterFlags()
}

var reloaders = make(map[string]Reloader)

func RegisterReloader(name string, reloader Reloader) {
	if reloader == nil {
		panic("reloader: Register reloader is nil")
	}

	if _, dup := reloaders[name]; dup {
		panic("reloader: Register called twice for reloader " + name)
	}
	reloaders[name] = reloader
}

func OpenReloader(reloaderName string) (Reloader, error) {
	reloader, ok := reloaders[reloaderName]
	if !ok {
		return nil, fmt.Errorf("reloader: unkown driver %q (forgotten import?)", reloaderName)
	}

	return reloader, nil
}

func RegisterReloaderFlags() {
	for _, reloader := range reloaders {
		reloader.RegisterFlags()
	}
}
