// Package downloader provides functionality to download a url and check its sha256 sum.
package downloader

import (
	"archive/zip"
	"bytes"
	sha256pkg "crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Downloader struct {
	filter func(path string) bool
}

func New(filter func(string) bool) *Downloader {
	return &Downloader{
		filter: filter,
	}
}

func (d *Downloader) Download(fs afero.Fs, destDir, url, sha256 string) error {
	rsp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http get")
	}
	defer rsp.Body.Close()
	logrus.Debugf("fetched %s", url)

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status: %d", rsp.StatusCode)
	}

	buf := bytes.NewBuffer([]byte{})
	h := sha256pkg.New()
	w := io.MultiWriter(buf, h)
	if _, err := io.Copy(w, rsp.Body); err != nil {
		return errors.Wrap(err, "sum body")
	}
	acSHA256 := hex.EncodeToString(h.Sum([]byte{}))
	if acSHA256 != sha256 {
		return fmt.Errorf("sum mismatch (%s != %s)", acSHA256, sha256)
	}
	logrus.Debugf("summed %s", url)

	if err := d.unzip(fs, destDir, buf); err != nil {
		return errors.Wrap(err, "unzip")
	}
	logrus.Debug("unzip success")

	return nil
}

func (d *Downloader) unzip(
	fs afero.Fs,
	destDir string,
	buf *bytes.Buffer,
) error {
	r := bytes.NewReader(buf.Bytes())
	size := int64(buf.Len())

	zipR, err := zip.NewReader(r, size)
	if err != nil {
		return errors.Wrap(err, "new reader")
	}

	for _, file := range zipR.File {
		logrus.Debugf("unzip: examining file %s", file.Name)

		fileR, err := file.Open()
		if err != nil {
			return errors.Wrap(err, "open file")
		}
		defer fileR.Close()

		path := filepath.Join(destDir, file.Name)
		if d.filter(path) {
			dir := filepath.Dir(path)
			logrus.Debugf("unzip: mkdir %s", dir)
			if err := fs.MkdirAll(dir, 0700); err != nil {
				return errors.Wrap(err, "mkdir")
			}

			logrus.Debugf("unzip: create %s", path)
			w, err := fs.Create(path)
			if err != nil {
				return errors.Wrap(err, "create file")
			}
			defer w.Close()

			if _, err := io.Copy(w, fileR); err != nil {
				return errors.Wrap(err, "copy")
			}
		}
	}

	return nil
}
