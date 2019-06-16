// Package deps provides dependency resolution.
package deps

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:generate counterfeiter . Downloader

type Downloader interface {
	Download(fs afero.Fs, destDir, url, sha256 string) error
}

type Deps struct {
	fs         afero.Fs
	cache      string
	downloader Downloader
}

func New(fs afero.Fs, cache string, downloader Downloader) *Deps {
	return &Deps{
		fs:         fs,
		cache:      cache,
		downloader: downloader,
	}
}

func (d *Deps) ResolveInclude(include string) (string, error) {
	dep := resolve(include)
	if dep == nil {
		return "", nil
	}
	logrus.Debugf("resolved include %s as dep %s", include, dep.name)

	if d.needsDownload(dep) {
		if err := d.downloadDep(dep); err != nil {
			return "", errors.Wrap(err, "download dep")
		}
	}

	for _, includePath := range dep.includePaths {
		path := filepath.Join(d.depPath(dep), includePath, include)
		logrus.Debugf("does %s exist", path)
		if exists, _ := afero.Exists(d.fs, path); exists {
			return path, nil
		}
	}

	return "", errors.New("could not find include in downloaded dependency")
}

func (d *Deps) ResolveSources(include string) ([]string, error) {
	dep := resolve(include)
	if dep == nil {
		return nil, nil
	}

	for _, source := range dep.sources {
		path := filepath.Join(d.depPath(dep), source)
		if exists, _ := afero.Exists(d.fs, path); !exists {
			if err := d.downloadDep(dep); err != nil {
				return nil, errors.Wrap(err, "download dep")
			}
		}
	}

	return dep.sources, nil
}

func (d *Deps) depPath(dep *dep) string {
	return filepath.Join(
		d.cache,
		"dependencies",
		dep.name,
	)
}

func (d *Deps) needsDownload(dep *dep) bool {
	if exists, _ := afero.Exists(d.fs, d.depPath(dep)); exists {
		return false
	} else {
		return true
	}
}

func (d *Deps) downloadDep(dep *dep) error {
	logrus.Infof("downloading dep %s", dep.name)

	if err := d.fs.MkdirAll(d.depPath(dep), 0700); err != nil {
		return errors.Wrap(err, "mkdir dep path")
	}

	if err := d.downloader.Download(
		d.fs,
		d.depPath(dep),
		dep.url,
		dep.sha256,
	); err != nil {
		return errors.Wrap(err, "download")
	}

	return nil
}

func resolve(include string) *dep {
	for _, dep := range deps {
		for _, header := range dep.headers {
			if include == header {
				return &dep
			}
		}
	}
	return nil
}
