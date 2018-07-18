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

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Specify input and output or -d for debug\n")
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

	definition := compose.Convert(graph)
	out, err := yaml.Marshal(definition)
	die(err)

	if Debug {
		fmt.Printf("Out:\n\n%s\n", out)
	} else {
		err := ioutil.WriteFile(os.Args[2], out, 0644)
		die(err)
	}
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
