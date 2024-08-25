package joker

import (
	"os/exec"

	"github.com/flosch/pongo2/v6"
	"github.com/toudi/joker/internal/jokerfile"
)

type Service struct {
	definition jokerfile.Service
	process    *exec.Cmd
}

func (s *Service) templateContext() *pongo2.Context {
	return &pongo2.Context{
		"service": map[string]interface{}{
			"data_dir": s.definition.Dir,
		},
	}
}
