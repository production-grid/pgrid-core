package main

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/security"
)

func main() {

	app := applications.Application{
		Name: "Production Grid Core Demo",
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.Start()

}
