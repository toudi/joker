package joker

import (
	"context"
	"fmt"

	"github.com/flosch/pongo2/v6"
	"github.com/toudi/joker/internal/jokerfile"
	"github.com/toudi/joker/internal/statefile"
)

type Joker struct {
	config            *jokerfile.Jokerfile
	services          []*Service
	state             *statefile.State
	env               pongo2.Context
	streamChan        chan StreamLine
	shutdownFunctions []func()
	ctx               context.Context
}

func Joker_init(ctx context.Context, configfile string) (*Joker, error) {
	config, err := jokerfile.Parse(configfile)
	if err != nil {
		return nil, err
	}

	jkr := &Joker{
		config:     config,
		streamChan: make(chan StreamLine),
		ctx:        ctx,
	}

	for _, serviceDefinition := range config.Services {
		jkr.services = append(jkr.services, &Service{definition: serviceDefinition})
	}

	jkr.env = make(pongo2.Context)

	if config.Environment != nil {
		jkr.env.Update(pongo2.Context{"env": config.Environment})
	}

	return jkr, nil
}

func (j *Joker) SetStatefile(path string) error {
	state, err := statefile.Parse(path)
	if err != nil {
		return err
	}
	j.state = state
	return nil
}

func (j *Joker) SaveState() error {
	return j.state.Save()
}

func (j *Joker) interpolateEnvVars(value string) string {
	pongoTmpl, err := pongo2.FromString(value)
	if err != nil {
		fmt.Printf("error preparing template: %v\n", err)
		return value
	}

	value, _ = pongoTmpl.Execute(j.env)
	return value
}
