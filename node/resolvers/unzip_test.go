package resolvers_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnzip(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name      string
		failure   bool
		dirSHA256 string
	}{
		{
			name:    "BadZipFile",
			failure: true,
		},
		{
			name:      "Success",
			failure:   false,
			dirSHA256: "a162e68f18a23ee6b09f896e95ae1de668eeb3ffff3d468980d03a63eb1504d7",
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			outputDir, err := ioutil.TempDir("", "btool_unzip_test")
			require.Nil(t, err)
			defer os.RemoveAll(outputDir)

			u := resolvers.NewUnzip(outputDir)
			zipPath := filepath.Join("testdata", datum.name, "zip.zip")
			zipN := node.New(zipPath)
			n := node.New("whatever").Dependency(zipN)
			err = u.Resolve(n)
			if datum.failure {
				assert.NotNil(t, err)
				return
			} else {
				require.Nil(t, err)
			}

			fs := afero.NewOsFs()
			dirSHA256, err := dirSHA256(fs, outputDir)
			require.Nil(t, err)
			assert.Equal(t, dirSHA256, datum.dirSHA256)
		})
	}
}

func dirSHA256(fs afero.Fs, dir string) (string, error) {
	h := sha256.New()

	if err := afero.Walk(fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "walk "+path)
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return errors.Wrap(err, "rel")
		}

		logrus.Debugf("summing path %s", relPath)
		h.Write([]byte(relPath))

		if info.IsDir() {
			return nil
		}

		logrus.Debugf("summing data for path %s", path)
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			return errors.Wrap(err, "read "+path)
		}

		h.Write(data)

		return nil
	}); err != nil {
		return "", errors.Wrap(err, "walk")
	}

	return hex.EncodeToString(h.Sum([]byte{})), nil
}
