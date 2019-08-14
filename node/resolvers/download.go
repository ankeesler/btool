package resolvers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type download struct {
	c      *http.Client
	url    string
	sha256 string
}

// NewDownload returns a node.Resolver that downloads a node.Node from an
// HTTP/HTTPS URL.
func NewDownload(
	c *http.Client,
	url string,
	sha256 string,
) node.Resolver {
	return &download{
		c:      c,
		url:    url,
		sha256: sha256,
	}
}

func (d *download) Resolve(n *node.Node) error {
	rsp, err := d.c.Get(d.url)
	if err != nil {
		return errors.Wrap(err, "http get")
	} else if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status: %d", rsp.StatusCode)
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "read body")
	}

	outputFile := n.Name
	logrus.Debugf(
		"got %d bytes from %s, writing to %s",
		len(data),
		d.url,
		outputFile,
	)

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}

	if err := ioutil.WriteFile(outputFile, data, 0644); err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}
