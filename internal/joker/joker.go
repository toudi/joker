package joker

import (
	"context"

	"github.com/flosch/pongo2/v6"
	"github.com/phuslu/log"
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
	// will only be launched if there's at least one service
	// that requests hot reloading
	hotReloadWatcher *HotReloadWatcher
	ctx              context.Context
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
		env:        make(pongo2.Context),
	}

	if config.Environment != nil {
		jkr.env.Update(pongo2.Context{"env": config.Environment})
	}

	for _, serviceDefinition := range config.Services {
		service := &Service{definition: serviceDefinition}
		if err = service.prepareDir(jkr); err != nil {
			return nil, err
		}
		jkr.services = append(jkr.services, service)
		if serviceDefinition.HotReload != nil {
			if err = jkr.prepareHotReloading(service); err != nil {
				return nil, err
			}
		}
	}

	if err = jkr.startHotReloadWatcher(); err != nil {
		return nil, err
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
		log.Error().Err(err).Msg("error preparing template")
		return value
	}

	value, _ = pongoTmpl.Execute(j.env)
	return value
}
