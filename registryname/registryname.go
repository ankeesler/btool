// Package registryname helps calculate the name for a running registry.
package registryname

import (
	"github.com/ankeesler/btool/log"
	"github.com/cloudfoundry-community/go-cfenv"
)

// Get returns the name for the registry running in the current process.
// It returns an error if a name cannot be gotten.
func Get(deefault string) (string, error) {
	app, err := cfenv.Current()
	if err != nil {
		log.Debugf("note: cannot get cfenv: %s", err.Error())
	} else if len(app.ApplicationURIs) == 0 {
		log.Debugf("note: no cfenv application uris")
	} else {
		return app.ApplicationURIs[0], nil
	}

	return deefault, nil
}
