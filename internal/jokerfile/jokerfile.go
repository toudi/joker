package jokerfile

type Service struct {
	Name      string
	Dir       string            `yaml:"dir"`
	Build     []string          `yaml:"build"`
	Command   interface{}       `yaml:"command"`
	Cleanup   interface{}       `yaml:"cleanup"`
	Depends   []string          `yaml:"depends"`
	Bootstrap interface{}       `yaml:"bootstrap"`
	HotReload interface{}       `yaml:"hot-reload"`
	Env       map[string]string `yaml:"env"`
}

type Jokerfile struct {
	Services    []Service
	Environment map[string]interface{} `yaml:"env"`
}
