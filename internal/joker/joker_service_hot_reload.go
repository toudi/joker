package joker

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/phuslu/log"
)

var (
	ErrServiceDirectoryNotDefined = errors.New("service directory not defined")
	ErrUnknownTypeForHotReload    = errors.New(
		"unrecognized type for hot-reload property. please use either bool or list of strings",
	)
	ErrTraversingDirectories  = errors.New("error traversing hot-reload tree")
	ErrProducingDirectoryName = errors.New("error figuring out the directory name")
)

func (s *Service) HotReloadHandler(joker *Joker) error {
	var srcDirs []string
	var directories []string

	// first, let's check what is the type of HotReload property:
	if _, ok := s.definition.HotReload.(bool); ok {
		// if it's a boolean then we simply have to add the service's directory to watchlist
		if s.definition.Dir == "" {
			return ErrServiceDirectoryNotDefined
		}
		srcDirs = append(srcDirs, s.definition.Dir)
	} else if srcDir, ok := s.definition.HotReload.(string); ok {
		srcDirs = append(srcDirs, srcDir)
	} else if srcDirsList, ok := s.definition.HotReload.([]string); ok {
		// if that's not the case, maybe it's a list of strings?
		srcDirs = srcDirsList
	} else {
		return ErrUnknownTypeForHotReload
	}

	// because fsnotify does not support subdirectories, let's list them manually

	for _, dir := range srcDirs {
		if err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if fullPath, err := filepath.Abs(filepath.Join(dir, path)); err != nil {
					return errors.Join(ErrProducingDirectoryName, err)
				} else {
					directories = append(directories, fullPath)
				}
			}
			return nil
		}); err != nil {
			return errors.Join(ErrTraversingDirectories, err)
		}
	}

	// great. now that we've got a list of directories to watch - let's pass that to fsnotify
	return s.executeHotReloadHandler(joker, directories)
}

func (s *Service) executeHotReloadHandler(joker *Joker, directories []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	joker.Defer(func() {
		watcher.Close()
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
					if err := s.Down(); err != nil {
						log.Error().Err(err).Msg("could not stop service")
					}
					if err = s.Up(joker.ctx, joker); err != nil {
						log.Error().Err(err).Msg("could not start service")
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	// Add paths.
	for _, dir := range directories {
		if err = watcher.Add(dir); err != nil {
			return err
		}
	}

	return nil
}
