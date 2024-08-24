package joker

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/phuslu/log"
)

const reloadThreshold = 500 * time.Millisecond

type HotReloadWatcher struct {
	// since there's a single watcher, but multiple services we have to
	// keep a reverse index from path to the service so that when there's
	// a change, we can quickly detect which service does it belong to
	// and restart it.
	pathsByService map[*Service][]string
	pathsToWatch   []string
	// this map keeps track of the modified files. we will use it against
	// known roots to deduce which services should be reloaded.
	detectedChanges map[string]bool
	// we keep track as to when the last change occured so we only restart
	// once (hopefully)
	lastChangeTimestamp time.Time
}

func (w *HotReloadWatcher) registerWatcher(service *Service) error {
	// first, let's check what is the type of HotReload property:
	if _, ok := service.definition.HotReload.(bool); ok {
		// if it's a boolean then we simply have to add the service's directory to watchlist
		if service.definition.Dir == "" {
			return ErrServiceDirectoryNotDefined
		}
		w.pathsByService[service] = []string{filepath.Clean(service.definition.Dir)}
	} else if srcDir, ok := service.definition.HotReload.(string); ok {
		w.pathsByService[service] = []string{filepath.Clean(srcDir)}
	} else if srcDirsList, ok := service.definition.HotReload.([]string); ok {
		// if that's not the case, maybe it's a list of strings?
		for _, srcDir := range srcDirsList {
			w.pathsByService[service] = append(w.pathsByService[service], filepath.Clean(srcDir))
		}
	} else {
		return ErrUnknownTypeForHotReload
	}
	return nil
}

func (w *HotReloadWatcher) reloadAffectedServices(joker *Joker) {
	var reloadServices []*Service
	var servicesSet map[*Service]bool = make(map[*Service]bool)

	for path := range w.detectedChanges {
		// we have to lookup the service based on affected path
		for service, srcDirs := range w.pathsByService {
			for _, dir := range srcDirs {
				if strings.HasPrefix(path, dir) && !servicesSet[service] {
					reloadServices = append(reloadServices, service)
					servicesSet[service] = true
					break
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
		if err := service.Up(joker.ctx, joker); err != nil {
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
