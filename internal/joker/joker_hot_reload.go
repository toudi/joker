package joker

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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
			pathsByService:  make(map[*Service][]string),
			detectedChanges: make(map[string]bool),
		}
	}

	return j.hotReloadWatcher.registerWatcher(service)
}

func (j *Joker) startHotReloadWatcher() error {
	if j.hotReloadWatcher != nil {
		// first, we have to traverse all the roots and calculate the final list of paths to
		// be watched.

		for _, srcDirs := range j.hotReloadWatcher.pathsByService {
			for _, dir := range srcDirs {
				if err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
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
