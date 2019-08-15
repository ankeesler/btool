package resolvers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type unzip struct {
	outputDir string
}

// NewUnzip returns a node.Resolver that unzips a dependency
func NewUnzip(outputDir string) node.Resolver {
	return &unzip{
		outputDir: outputDir,
	}
}

func (u *unzip) Resolve(n *node.Node) error {
	if ex, ac := 1, len(n.Dependencies); ex != ac {
		return fmt.Errorf("expected %d dependency, got %d", ex, ac)
	}

	zipData, err := ioutil.ReadFile(n.Dependencies[0].Name)
	if err != nil {
		return errors.Wrap(err, "read file")
	}
	zipBuf := bytes.NewBuffer(zipData)
	logrus.Debugf("zip buf len = %d", zipBuf.Len())

	fs := afero.NewOsFs()
	if err := u.unzip(fs, u.outputDir, zipBuf); err != nil {
		return errors.Wrap(err, "unzip")
	}

	return nil
}

func (u *unzip) unzip(
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
		logrus.Debugf("unzip file %s", file.Name)
		if err := u.unzipFile(fs, destDir, file); err != nil {
			return errors.Wrap(err, "unzip file")
		}
	}

	return nil
}

func (u *unzip) unzipFile(
	fs afero.Fs,
	destDir string,
	file *zip.File,
) error {
	fileR, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer fileR.Close()

	path := filepath.Join(destDir, file.Name)
	if file.FileInfo().IsDir() {
		logrus.Debugf("unzip: mkdir %s", path)
		if err := fs.Mkdir(path, 0755); err != nil {
			return errors.Wrap(err, "mkdir")
		}
	} else {
		logrus.Debugf("unzip: create %s", path)
		w, err := fs.Create(path)
		if err != nil {
			return errors.Wrap(err, "create")
		}
		defer w.Close()

		if _, err := io.Copy(w, fileR); err != nil {
			return errors.Wrap(err, "copy")
		}
	}

	return nil
}
