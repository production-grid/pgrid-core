package loaders

import (
	"errors"
	"io"
	"os"

	"github.com/production-grid/pgrid-core/pkg/logging"
)

// EnvResourcePath is the default environment variable for the resource path
const EnvResourcePath = "PG_RC_HOME"

// FileResourceLoader loads resources relative to a give directory.
type FileResourceLoader struct {
	EnvironmentVariable string
	BasePath            string
}

// Reader returns a reader for the given path
func (loader *FileResourceLoader) Reader(path string) (io.Reader, error) {

	if loader.BasePath == "" {
		if loader.EnvironmentVariable == "" {
			loader.EnvironmentVariable = EnvResourcePath
		}
		loader.BasePath = os.Getenv(loader.EnvironmentVariable)
	}

	if loader.BasePath == "" {
		return nil, errors.New("no base path configured for file resource loader")
	}

	logging.Tracef("Loading File Resource: %v\n", path)
	return os.Open(loader.BasePath + "/" + path)

}

// Bytes returns a byte slice for a given path.
func (loader *FileResourceLoader) Bytes(path string) ([]byte, error) {
	return Bytes(loader, path)
}

// String returns a string for a given path.
func (loader *FileResourceLoader) String(path string) (string, error) {
	return String(loader, path)
}
