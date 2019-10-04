package main

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/security"
)

func main() {

	loader := &loaders.FileResourceLoader{}

	coreConfig, err := config.LoadCore(loader, "demo-config.yml")

	if err != nil {
		panic(err)
	}

	app := applications.Application{
		CoreConfiguration: *coreConfig,
		Name:              "Production Grid Core Demo",
		ConfigLoader:      loader,
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.PreMigrate()
	app.Start()
	app.PostMigrate()

}
