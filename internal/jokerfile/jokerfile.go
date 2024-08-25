package jokerfile

type Service struct {
	Name            string
	Dir             string            `yaml:"dir"`
	Build           []string          `yaml:"build"`
	Command         interface{}       `yaml:"command"`
	Cleanup         interface{}       `yaml:"cleanup"`
	Depends         []string          `yaml:"depends"`
	Bootstrap       interface{}       `yaml:"bootstrap"`
	HotReload       interface{}       `yaml:"hot-reload"`
	Env             map[string]string `yaml:"env"`
	KillSignal      interface{}       `yaml:"shutdown-signal"`
	NoLivenessCheck bool              `yaml:"no-liveness-check"`
}

type Jokerfile struct {
	DataDir     string
	Services    []Service
	Commands    map[string]interface{} `yaml:"commands"`
	Environment map[string]interface{} `yaml:"env"`
}
