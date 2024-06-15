package joker

import (
	"os/exec"

	"github.com/toudi/joker/internal/jokerfile"
)

type Service struct {
	definition jokerfile.Service
	process    *exec.Cmd
}
