package joker

import (
	"errors"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/phuslu/log"
	"github.com/samber/lo"
	"github.com/toudi/joker/internal/utils"
)

const reloadThreshold = 500 * time.Millisecond

var ErrNotASliceOfStrings = errors.New(
	"you've defined hot-reload as an array but it does not seem to be an array of strings",
)

type HotReloadConfig struct {
	pathsWithPatterns []string
	// this will be used to check if the modified file matches the defined pattern
}

type HotReloadWatcher struct {
	// since there's a single watcher, but multiple services we have to
	// keep a reverse index from path to the service so that when there's
	// a change, we can quickly detect which service does it belong to
	// and restart it.
	configByService map[*Service]HotReloadConfig
	pathsToWatch    []string
	// this map keeps track of the modified files. we will use it against
	// known roots to deduce which services should be reloaded.
	detectedChanges map[string]bool
	// we keep track as to when the last change occured so we only restart
	// once (hopefully)
	lastChangeTimestamp time.Time
	jokerWorkingDir     string
}

func (w *HotReloadWatcher) registerWatcher(service *Service) error {
	if service.definition.Dir == "" {
		return ErrServiceDirectoryNotDefined
	}

	// let's check if this is a slice of strings
	var err error

	// let's check what type of config we have
	if _, ok := service.definition.HotReload.(bool); ok {
		// somebody just defined:
		// hot-reload: true
		w.configByService[service] = HotReloadConfig{
			pathsWithPatterns: []string{
				w.absolutePathWithPattern(
					"*",
					service,
				),
			},
		}
	} else if srcDir, ok := service.definition.HotReload.(string); ok {
		// somebody just defined:
		// hot-reload: some/path
		w.configByService[service] = HotReloadConfig{
			pathsWithPatterns: []string{w.absolutePathWithPattern(filepath.Clean(srcDir), service)},
		}
	} else if srcDirsList, ok := service.definition.HotReload.([]interface{}); ok {
		// example:
		// hot-reload:
		//  - /a/path
		//  - /some/other/path
		// if that's not the case, maybe it's a list of strings?
		w.configByService[service] = HotReloadConfig{
			pathsWithPatterns: lo.Map(srcDirsList, func(srcDir interface{}, _ int) string {
				if tmpString, ok := srcDir.(string); !ok {
					err = ErrNotASliceOfStrings
					return ""
				} else {
					return w.absolutePathWithPattern(tmpString, service)
				}
			}),
		}
	} else {
		err = ErrUnknownTypeForHotReload
	}

	return err
}

func (w *HotReloadWatcher) reloadAffectedServices(joker *Joker) {
	var reloadServices []*Service
	var servicesSet map[*Service]bool = make(map[*Service]bool)

	for filePath := range w.detectedChanges {
		// we have to lookup the service based on affected path
		for service, config := range w.configByService {
			for _, pattern := range config.pathsWithPatterns {
				log.Trace().Str("path", filePath).Str("expected prefix", path.Dir(pattern)).Msg("")
				if strings.HasPrefix(filePath, path.Dir(pattern)) && !servicesSet[service] {
					// let's check if the pattern matches:
					patternMatch, err := filepath.Match(pattern, filePath)
					if err != nil {
						log.Error().
							Err(err).
							Msg("cannot detect if file name matches expected pattern")
					}
					if !patternMatch {
						log.Trace().
							Msgf("%s does not match a pattern %s", filePath, path.Base(pattern))
						continue
					}
					reloadServices = append(reloadServices, service)
					servicesSet[service] = true
					break
				} else {
					log.Trace().Msgf("%s either does not begin with %s or it is an unknown service", filePath, path.Dir(pattern))
				}
			}
		}
	}

	for key := range w.detectedChanges {
		delete(w.detectedChanges, key)
	}

	for _, service := range reloadServices {
		log.Trace().Str("service", service.definition.Name).Msg("reloading")
		if err := service.Down(service.shutdownOptions(true)); err != nil {
			log.Error().Err(err).Msg("could not stop service")
		}
		if err := service.Up(joker.ctx, joker, serviceStartOptions{}); err != nil {
			log.Error().Err(err).Msg("could not start service")
		}
	}
}

func (w *HotReloadWatcher) startWorker(joker *Joker) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	var ticker = time.NewTicker(reloadThreshold)

	joker.Defer(func() {
		watcher.Close()
		ticker.Stop()
	})

	var jokerRuntimeFiles = []string{
		filepath.Join(w.jokerWorkingDir, joker.initOptions.Jokerfile),
		filepath.Join(w.jokerWorkingDir, joker.initOptions.StateFile),
		filepath.Join(w.jokerWorkingDir, rpcListenerFile),
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Trace().Str("path", event.Name).Msg("detected filesystem change")
					if slices.Contains(jokerRuntimeFiles, event.Name) {
						log.Trace().Str("path", event.Name).Msg("detected runtime file; no-op")
						continue
					}
					// do not reload the service right away as there might be multiple
					// changes. Simply set the indicator flag and let the timer handle
					// the rest.
					w.detectedChanges[event.Name] = true
					w.lastChangeTimestamp = time.Now()
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			case now := <-ticker.C:
				if now.Sub(w.lastChangeTimestamp) >= reloadThreshold {
					w.reloadAffectedServices(joker)
				}
			}
		}
	}()

	// Add paths.
	for _, dir := range w.pathsToWatch {
		if err = watcher.Add(dir); err != nil {
			return err
		}
	}

	return nil
}

func (w *HotReloadWatcher) absolutePathWithPattern(input string, service *Service) string {
	path := utils.PathToPathWithPattern(input)

	if strings.HasPrefix(input, "./") {
		// this is a path that originates in joker's working dir
		return filepath.Clean(filepath.Join(w.jokerWorkingDir, path))
	} else if strings.HasPrefix(path, "/") {
		// this is an absolute path already
		return path
	}
	// otherwise it must be a relative path so let's make it absolute with relation
	// to service dir
	return filepath.Join(service.definition.Dir, path)
}
