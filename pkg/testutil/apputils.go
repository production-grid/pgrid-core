package testutil

import (
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/cache"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/production-grid/pgrid-core/pkg/security"

	"github.com/gorilla/handlers"
)

//TestServer provides a reference to the http test server
var TestServer *httptest.Server

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

	defaultUser := applications.DefaultAdminUser{
		FirstName: "System",
		LastName:  "Administrator",
		EMail:     "devops@productiongrid.com",
		Password:  "test123",
	}

	localCache := cache.LocalCache{}

	app := applications.Application{
		CoreConfiguration:   *coreConfig,
		DefaultAdminUser:    &defaultUser,
		SessionStore:        &localCache,
		GeneralPurposeCache: &localCache,
		Name:                "Production Grid Integration Test Application",
		ConfigLoader:        loader,
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.PreMigrate()
	app.PostMigrate()
	app.Start()

	appSingleton = &app

	// Print server logs to stdout
	handler := handlers.LoggingHandler(os.Stdout, app.Router)

	logging.Infoln("Starting HTTP Server...")
	TestServer = httptest.NewServer(handler)

}
