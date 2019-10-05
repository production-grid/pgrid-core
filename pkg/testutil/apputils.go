package demo

import (
	"sync"
	"testing"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/security"
)

var initOnce sync.Once

var appSingleton *applications.Application

// StartTestApplication starts the test application in order to support integration tests.
func StartTestApplication(t *testing.T) *applications.Application {

	initOnce.Do(func() {

		doTestStartup(t)

	})

	return appSingleton

}

func doTestStartup(t *testing.T) {

	loader := &loaders.FileResourceLoader{}

	coreConfig, err := config.LoadCore(loader, "demo-config.yml")

	if err != nil {
		panic(err)
	}

	app := applications.Application{
		CoreConfiguration: *coreConfig,
		Name:              "Production Grid Integration Test Application",
		ConfigLoader:      loader,
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.PreMigrate()
	app.PostMigrate()
	go app.Start()

	appSingleton = &app

}
