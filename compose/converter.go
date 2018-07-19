package compose

import (
	types "github.com/tariel-x/anzer/funcs"
	"strings"
)

const DefaultEnvOut = "OUT"

func Convert(graph types.SystemGraph) Compose {
	compose := base()

	for _, set := range graph.Services {
		addition := addServices(set)
		for name, def := range addition {
			compose.Services[name] = def
		}
	}

	return compose
}

func addServices(set types.ServiceSet) map[string]Service {
	addition := map[string]Service{}
	for _, source := range set {
		name, service := makeService(source)
		addition[name] = service
	}
	return addition
}

func makeService(source types.Service) (string, Service) {
	uniqueName := strings.Replace(source.UniqueName, ".", "_", -1)
	envs := source.Config.Envs
	if envs == nil {
		envs = map[string]string{}
	}
	if source.Config.EnvOut != "" {
		envs[source.Config.EnvOut] = uniqueName
	} else {
		envs[DefaultEnvOut] = uniqueName
	}
	return uniqueName, Service{
		Image:       source.Name,
		Container:   uniqueName,
		Environment: envs,
	}
}

func base() Compose {
	return Compose{
		Version:  "3",
		Services: map[string]Service{},
	}
}
