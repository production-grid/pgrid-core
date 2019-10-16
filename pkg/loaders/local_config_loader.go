package loaders

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/production-grid/pgrid-core/pkg/logging"
)

// LocalConfigLoader loads resources relative to a give directory.
type LocalConfigLoader struct {
	configHome string
}

// Reader returns a reader for the given path
func (loader *LocalConfigLoader) Reader(path string) (io.Reader, error) {

	if loader.configHome == "" {
		var configHome string

		if runtime.GOOS == "windows" {
			configHome = os.Getenv("userprofile")
		} else {
			configHome = os.Getenv("XDG_CONFIG_HOME")
			if configHome == "" {
				user, err := user.Current()
				if err != nil {
					return nil, err
				}
				configHome = user.HomeDir + "/.config"
			}
		}
		loader.configHome = configHome
	}
	path = filepath.Join(loader.configHome, path)
	logging.Tracef("Loading Config Resource: %v\n", path)
	return os.Open(path)

}

// Bytes returns a byte slice for a given path.
func (loader *LocalConfigLoader) Bytes(path string) ([]byte, error) {
	return Bytes(loader, path)
}

// String returns a string for a given path.
func (loader *LocalConfigLoader) String(path string) (string, error) {
	return String(loader, path)
}
