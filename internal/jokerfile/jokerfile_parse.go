package jokerfile

import (
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

type parsableJokerfile struct {
	Environment map[string]interface{} `yaml:"env"`
	Services    map[string]Service     `yaml:"services"`
	Commands    map[string]interface{} `yaml:"commands"`
}

func Parse(filePath string) (*Jokerfile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var srcJokerfile parsableJokerfile
	var jkrfile Jokerfile
	if err = yaml.NewDecoder(file).Decode(&srcJokerfile); err != nil {
		return nil, err
	}

	jkrfile.Environment = srcJokerfile.Environment
	jkrfile.Commands = srcJokerfile.Commands

	for serviceName, service := range srcJokerfile.Services {
		service.Name = serviceName
		jkrfile.Services = append(jkrfile.Services, service)
	}

	// now let's sort the services according to their dependencies.
	slices.SortFunc(jkrfile.Services, func(a, b Service) int {
		var dependency string
		if a.Depends == nil && b.Depends == nil {
			return 0
		}
		if b.Depends != nil {
			for _, dependency = range b.Depends {
				if a.Name == dependency {
					return -1
				}
			}
		}
		if a.Depends != nil {
			for _, dependency = range a.Depends {
				if b.Name == dependency {
					return 1
				}
			}
		}
		return 0
	})

	return &jkrfile, nil
}
