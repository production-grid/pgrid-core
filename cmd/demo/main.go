package main

import (
	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/cache"
	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/security"
)

func main() {

	rcLoader := &loaders.FileResourceLoader{}
	configLoader := &loaders.LocalConfigLoader{}

	coreConfig, err := config.LoadCore(configLoader, "pgrid/pgrid_dev.yml")

	if err != nil {
		panic(err)
	}

	lingo := applications.TenantLingo{
		TenantSingularKey:     "arts.organization",
		TenantPluralKey:       "arts.organizations",
		TenantSingularDefault: "Arts Organization",
		TenantPluralDefault:   "Arts Organizations",
		Types: []applications.TenantType{
			applications.TenantType{
				SingularKey:     "theatre.company",
				PluralKey:       "theatre.companies",
				SingularDefault: "Theatre Company",
				PluralDefault:   "Theatre Companies",
			},
			applications.TenantType{
				SingularKey:     "school.drama.program",
				PluralKey:       "school.drama.programs",
				SingularDefault: "School Drama Program",
				PluralDefault:   "School Drama Programs",
			},
			applications.TenantType{
				SingularKey:     "school.music.program",
				PluralKey:       "school.music.programs",
				SingularDefault: "School Music Program",
				PluralDefault:   "School Music Programs",
			},
			applications.TenantType{
				SingularKey:     "dance.company",
				PluralKey:       "dance.companies",
				SingularDefault: "Dance Company",
				PluralDefault:   "Dance Companies",
			},
			applications.TenantType{
				SingularKey:     "opera.company",
				PluralKey:       "opera.companies",
				SingularDefault: "Opera Company",
				PluralDefault:   "Opera Companies",
			},
			applications.TenantType{
				SingularKey:     "choir",
				PluralKey:       "choirs",
				SingularDefault: "Choir",
				PluralDefault:   "Choirs",
			},
			applications.TenantType{
				SingularKey:     "symphony",
				PluralKey:       "symphonies",
				SingularDefault: "Symphony",
				PluralDefault:   "Symphonies",
			},
			applications.TenantType{
				SingularKey:     "other.arts.group",
				PluralKey:       "other.arts.groups",
				SingularDefault: "Other Arts Group",
				PluralDefault:   "Other Arts Groups",
			},
		},
	}

	localCache := cache.LocalCache{}

	defaultUser := applications.DefaultAdminUser{
		FirstName: "System",
		LastName:  "Administrator",
		EMail:     "devops@productiongrid.com",
		Password:  "test123",
	}

	app := applications.Application{
		CoreConfiguration: *coreConfig,
		TenantLingo:       &lingo,
		SessionStore:      &localCache,
		Cache:             &localCache,
		DefaultAdminUser:  &defaultUser,
		Name:              "Production Grid Core Demo",
		ConfigLoader:      rcLoader,
		Modules: []applications.FeatureModule{
			&security.Module{},
		},
	}

	app.PreMigrate()
	app.PostMigrate()
	app.StartAndListen()

}
