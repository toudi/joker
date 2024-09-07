package joker

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

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

func (j *Joker) prepareHotReloading(service *Service) error {
	if j.hotReloadWatcher == nil {
		j.hotReloadWatcher = &HotReloadWatcher{
			configByService: make(map[*Service]HotReloadConfig),
			detectedChanges: make(map[string]bool),
			jokerWorkingDir: j.initOptions.Workdir,
		}
	}

	return j.hotReloadWatcher.registerWatcher(service)
}

func (j *Joker) startHotReloadWatcher() error {
	if j.hotReloadWatcher != nil {
		// first, we have to traverse all the roots and calculate the final list of paths to
		// be watched.

		for _, config := range j.hotReloadWatcher.configByService {
			for _, pattern := range config.pathsWithPatterns {
				dir := filepath.Dir(pattern)
				if err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							// the path does not exist which is not a fatal error, we can just no-op.
							log.Trace().Str("path", filepath.Clean(filepath.Join(dir, path))).Msg("does not seem to exist; no-op")
							return nil
						}
					}
					if d.IsDir() {
						if fullPath, err := filepath.Abs(filepath.Join(dir, path)); err != nil {
							return errors.Join(ErrProducingDirectoryName, err)
						} else {
							j.hotReloadWatcher.pathsToWatch = append(j.hotReloadWatcher.pathsToWatch, fullPath)
						}
					}
					return nil
				}); err != nil {
					return errors.Join(ErrTraversingDirectories, err)
				}
			}
		}

		// if we are here then it means that we can launch the watcher.
		return j.hotReloadWatcher.startWorker(j)
	}

	return nil
}
