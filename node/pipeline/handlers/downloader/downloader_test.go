package downloader_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/pipeline/handlers/downloader"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestDownload(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name string

		url    string
		sha256 string

		failure   bool
		dirSHA256 string
	}{
		{
			name: "BadAddress",

			url: "bad-address",

			failure: true,
		},
		{
			name: "BadHttpResponse",

			url: "https://github.com/this/path/does/not/exist",

			failure: true,
		},
		{
			name: "BadSHA256",

			url:    "https://github.com/ankeesler/anwork/archive/v9.zip",
			sha256: "788a4047f7ac2518508bb3080ff1f0ef196a59a09eb3938b3e41d48cbf9e64de",

			failure: true,
		},
		{
			name: "BadZipFile",

			url:    "https://github.com/ankeesler/anwork/releases/download/v9/v9_anwork_darwin_amd64",
			sha256: "35b7b6a0360ae1801ccd0adb89f6de6a0a8c08dd168337cf5ac7649d6e81cdca",

			failure: true,
		},
		{
			name: "Success",

			url:    "https://github.com/ankeesler/anwork/archive/v9.zip",
			sha256: "1969d4cef052d3365f9c7d53467178e84f2b596a59717566def9d1b7cfe40f8c",

			failure:   false,
			dirSHA256: "3a232593e95328866dc36a1f565ec027e6aeedabc9a44e80d12307acd1b4c7e8",
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			destDir := "/some/dest/dir"
			if err := downloader.New(func(path string) bool {
				return strings.HasSuffix(path, ".go")
			}).Download(
				fs,
				destDir,
				datum.url,
				datum.sha256,
			); err != nil {
				if !datum.failure {
					t.Fatal(err)
				}
				return
			} else if datum.failure {
				t.Fatal("expected failure")
			}

			dirSHA256, err := dirSHA256(fs, destDir)
			if err != nil {
				t.Fatal(err)
			}

			if dirSHA256 != datum.dirSHA256 {
				t.Fatalf("expected %s, actual %s", datum.dirSHA256, dirSHA256)
			}
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
