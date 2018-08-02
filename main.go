package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tariel-x/anzer-compose/compose"

	types "github.com/tariel-x/anzer/funcs"

	"gopkg.in/yaml.v2"
)

var (
	Debug = false
)

const (
	ComposeBase            = "base.yml"
	ComposeMQContainer     = "rmq"
	ComposeMQConnectionEnv = "RMQ"
	ComposeMQConnection    = "amqp://guest:guest@rmq:5672/"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Specify input, output or -d for debug and registry url\n")
		return
	}

	if os.Args[2] == "-d" {
		Debug = true
	}

	in, err := ioutil.ReadFile(os.Args[1])
	die(err)

	graph := types.SystemGraph{}
	err = json.Unmarshal(in, &graph)
	die(err)

	base, err := loadBase(ComposeBase)
	die(err)

	conf := compose.Config{
		Base:            base,
		MQContainer:     ComposeMQContainer,
		MQConnection:    ComposeMQConnection,
		MQConnectionEnv: ComposeMQConnectionEnv,
	}
	if len(os.Args) == 4 {
		conf.Registry = os.Args[3]
	}

	definition := compose.Convert(graph, conf)
	out, err := yaml.Marshal(definition)
	die(err)

	if Debug {
		fmt.Printf("Out:\n\n%s\n", out)
	} else {
		err := ioutil.WriteFile(os.Args[2], out, 0644)
		die(err)
	}
}

func loadBase(name string) (compose.Compose, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return compose.Compose{}, err
	}
	base := compose.Compose{}
	err = yaml.Unmarshal(data, &base)
	return base, err
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
