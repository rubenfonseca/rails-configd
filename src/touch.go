package src

import (
	"flag"
	"log"
	"os"
)

type TouchReloader struct {
	TouchFile *string
}

func (reloader *TouchReloader) Reload() {
	log.Printf("[TOUCH RELOADER] Touching %s", *reloader.TouchFile)

	file, err := os.OpenFile(*reloader.TouchFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	file.Truncate(0)
}

func (reloader *TouchReloader) RegisterFlags() {
	reloader.TouchFile = flag.String("touch-file", "tmp/restart.txt", "The file to touch when we need to reload")
}

func init() {
	touchReloader := TouchReloader{}
	RegisterReloader("touch", &touchReloader)
}
