package compose

type Compose struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks *struct {
		Frontend interface{} `yaml:"frontend"`
		Backend  interface{} `yaml:"backend"`
	} `yaml:"networks,omitempty"`
	Volumes *struct {
		DbData interface{} `yaml:"db-data"`
	} `yaml:"volumes,omitempty"`
}

type Service struct {
	Image       string            `yaml:"image"`
	Container   string            `yaml:"container"`
	Ports       []string          `yaml:"ports,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Deploy      *struct {
		Replicas     int `yaml:"replicas"`
		UpdateConfig struct {
			Parallelism int    `yaml:"parallelism"`
			Delay       string `yaml:"delay"`
		} `yaml:"update_config"`
		RestartPolicy struct {
			Condition string `yaml:"condition"`
		} `yaml:"restart_policy"`
	} `yaml:"deploy,omitempty"`
}
