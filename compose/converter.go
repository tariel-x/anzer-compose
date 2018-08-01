package compose

import (
	types "github.com/tariel-x/anzer/funcs"
	"strings"
)

const (
	DefaultEnvOut          = "OUT"
	DefaultEnvIn           = "IN"
	DefaultEnvProdIn1      = "IN1"
	DefaultEnvProdIn2      = "IN2"
	DefaultProductionImage = "compose_service"
	DefaultEnvType1Name    = "TYPE1"
	DefaultEnvType2Name    = "TYPE2"
)

type Config struct {
	Registry        string
	Base            Compose
	MQContainer     string
	MQConnection    string
	MQConnectionEnv string
}

func Convert(graph types.SystemGraph, conf Config) Compose {
	compose := conf.Base

	for _, set := range graph.Services {
		addition := addServices(set, conf)
		for name, def := range addition {
			compose.Services[name] = def
		}
	}

	for _, dependency := range graph.Dependencies {
		compose = setInput(dependency, graph.Services, compose)
	}

	return compose
}

func addServices(set types.ServiceSet, conf Config) map[string]Service {
	addition := map[string]Service{}
	for _, source := range set {
		var name string
		var service Service
		if source.Type == types.TypeProduction {
			name, service = makeProductionService(source, conf)
		} else {
			name, service = makeService(source, conf)
		}
		addition[name] = service
	}
	return addition
}

func makeService(source types.Service, conf Config) (string, Service) {
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
	envs[conf.MQConnectionEnv] = conf.MQConnection
	return uniqueName, Service{
		Image:       conf.Registry + source.Name,
		Container:   uniqueName,
		Environment: envs,
		DependsOn:   []string{conf.MQContainer},
	}
}

func makeProductionService(source types.Service, conf Config) (string, Service) {
	name, service := makeService(source, conf)

	service.Image = DefaultProductionImage
	if len(source.ProductionTypes) > 1 {
		service.Environment[DefaultEnvType1Name] = source.ProductionTypes[0]
		service.Environment[DefaultEnvType2Name] = source.ProductionTypes[1]
	}

	return name, service
}

func setInput(dependency types.Dependency, sources types.Services, compose Compose) Compose {
	original := strings.Split(dependency.To, ".")
	originalFrom := original[0]
	sourceEnvOut := sources[originalFrom][0].Config.EnvOut
	if sourceEnvOut == "" {
		sourceEnvOut = DefaultEnvOut
	}

	original = strings.Split(dependency.To, ".")
	originalTo := original[0]
	sourceEnvsIn := []string{}
	if sources[originalTo][0].Type == types.TypeProduction {
		sourceEnvsIn = []string{
			DefaultEnvProdIn1,
			DefaultEnvProdIn2,
		}
	} else {
		sourceEnvsIn = []string{sources[originalTo][0].Config.EnvIn}
		if len(sourceEnvsIn) == 0 {
			sourceEnvsIn = []string{DefaultEnvIn}
		}
	}

	fromName := strings.Replace(dependency.From, ".", "_", -1)
	outAddress := compose.Services[fromName].Environment[sourceEnvOut]

	toName := strings.Replace(dependency.To, ".", "_", -1)
	for _, sourceEnvIn := range sourceEnvsIn {
		in, ok := compose.Services[toName].Environment[sourceEnvIn]
		if !ok || in == "" {
			compose.Services[toName].Environment[sourceEnvIn] = outAddress
		}
	}
	return compose
}
