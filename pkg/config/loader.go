package config

import (
	"github.com/production-grid/pgrid-core/pkg/loaders"

	yaml "gopkg.in/yaml.v2"
)

// LoadCore the given path as CoreConfiguration
func LoadCore(loader loaders.ResourceLoader, path string) (*CoreConfiguration, error) {

	coreConfig := CoreConfiguration{}
	err := Load(loader, path, &coreConfig)

	if err != nil {
		return nil, err
	}

	return &coreConfig, nil

}

// Load loads configuration as yaml and attempts to populate the configTarget interface.
// This might get used directly if an application extends the configuration.
func Load(loader loaders.ResourceLoader, path string, configTarget interface{}) error {

	content, err := loader.Bytes(path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, configTarget)

}
