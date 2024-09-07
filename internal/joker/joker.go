package joker

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"strings"

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
	initOptions      InitOptions
}

var (
	errDataDirNotSpecified  = errors.New("data_dir not set in jokerfile")
	errUnableToSetStatefile = errors.New("unable to set statefile")
)

type InitOptions struct {
	Workdir   string
	Jokerfile string
	StateFile string
}

func Joker_init(ctx context.Context, options InitOptions) (*Joker, error) {
	config, err := jokerfile.Parse(options.Jokerfile)
	if err != nil {
		return nil, err
	}

	jkr := &Joker{
		config:      config,
		streamChan:  make(chan StreamLine),
		ctx:         ctx,
		env:         make(pongo2.Context),
		initOptions: options,
	}

	if jkr.initOptions.Workdir, err = filepath.Abs(path.Dir(jkr.initOptions.Jokerfile)); err != nil {
		return nil, err
	}

	if config.Environment != nil {
		jkr.env.Update(pongo2.Context{"env": config.Environment})
		interpolatedEnv := jkr.interpolateRecursively(
			config.Environment,
		)
		jkr.env.Update(pongo2.Context{"env": interpolatedEnv})
	}

	jkr.interpolateRecursively(jkr.env)

	var servicesEnv = map[string]interface{}{}

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
		servicesEnv[strings.ReplaceAll(serviceDefinition.Name, "-", "_")] = map[string]interface{}{
			"data_dir": service.definition.Dir,
		}
	}

	jkr.env.Update(pongo2.Context(map[string]interface{}{"services": servicesEnv}))
	if err = jkr.SetStatefile(options.StateFile); err != nil {
		return nil, errors.Join(errUnableToSetStatefile, err)
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

func (j *Joker) GetDataDir() (string, error) {
	if j.config.DataDir != "" {
		return j.config.DataDir, nil
	}
	return "", errDataDirNotSpecified
}

func (j *Joker) interpolateEnvVars(value string, additionalEnv *pongo2.Context) string {
	pongoTmpl, err := pongo2.FromString(value)
	if err != nil {
		log.Error().Err(err).Msg("error preparing template")
		return value
	}

	var context = j.env
	if additionalEnv != nil {
		context.Update(*additionalEnv)
	}

	value, err = pongoTmpl.Execute(context)
	if err != nil {
		log.Error().Err(err).Msg("error interpolating string")
	}
	return value
}

func (j *Joker) interpolateRecursively(env interface{}) interface{} {
	if aMap, ok := env.(map[string]interface{}); ok {
		for key, value := range aMap {
			aMap[key] = j.interpolateRecursively(value)
		}
		return aMap
	} else if aList, ok := env.([]interface{}); ok {
		for idx, item := range aList {
			aList[idx] = j.interpolateRecursively(item)
		}
		return aList
	} else if value, ok := env.(string); ok {
		return j.interpolateEnvVars(value, nil)
	}
	return env
}
