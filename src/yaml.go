package src

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/v1/yaml"
)

type YamlRenderer struct {
	YamlFile *string
}

func (renderer *YamlRenderer) Render(env Env) {
	log.Printf("[YAML RENDERER] Rendering to %s", *renderer.YamlFile)

	out, err := yaml.Marshal(env.Data)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(*renderer.YamlFile, out, 0644)
	if err != nil {
		panic(err)
	}
}

func (renderer *YamlRenderer) RegisterFlags() {
	renderer.YamlFile = flag.String("yaml-file", "config/config.yml", "The output of the Yaml file")
}

func init() {
	yamlRenderer := YamlRenderer{}
	RegisterRenderer("yaml", &yamlRenderer)
}
