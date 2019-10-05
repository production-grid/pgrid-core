package itests

import (
	"strings"
	"testing"

	"github.com/production-grid/pgrid-core/pkg/config"
	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/ids"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestSchemaDeleteProtection(t *testing.T) {

	assert := assert.New(t)

	loader := &loaders.FileResourceLoader{
		BasePath: "testdata",
	}

	coreConfig, err := config.LoadCore(loader, "schema-config.yml")

	assert.NoError(err)

	baseSchema := coreConfig.DatabaseConfiguration.Primary.Schema

	err = relational.Init(coreConfig.DatabaseConfiguration)

	assert.NoError(err)

	freshStart(t, coreConfig)

	testMigrations(t, coreConfig, loader, []string{"create-table-schema.json"})

	err = relational.PreMigrate(loader, coreConfig.DatabaseConfiguration, []string{"drop-column-schema.json"})

	assert.Error(err, "failure expected when attempting to drop columns")

	cleanUp(t, coreConfig, baseSchema)

}

// TestSchemaMigration tests the schema migration system.
func TestSchemaMigration(t *testing.T) {

	assert := assert.New(t)

	loader := &loaders.FileResourceLoader{
		BasePath: "testdata",
	}

	coreConfig, err := config.LoadCore(loader, "schema-config.yml")

	assert.NoError(err)

	baseSchema := coreConfig.DatabaseConfiguration.Primary.Schema

	err = relational.Init(coreConfig.DatabaseConfiguration)

	assert.NoError(err)

	freshStart(t, coreConfig)

	testMigrations(t, coreConfig, loader, []string{"create-table-schema.json"})
	testMigrations(t, coreConfig, loader, []string{"add-columns-schema.json"})
	testMigrations(t, coreConfig, loader, []string{"modify-columns-schema.json"})
	testMigrations(t, coreConfig, loader, []string{"add-fk-schema.json"})
	testMigrations(t, coreConfig, loader, []string{"add-index-schema.json"})

	cleanUp(t, coreConfig, baseSchema)
}

//cleanUp deletes the temporary database
func cleanUp(t *testing.T, coreConfig *config.CoreConfiguration, baseSchema string) {

	assert := assert.New(t)

	tempSchema := coreConfig.DatabaseConfiguration.Primary.Schema

	err := relational.Replica.Close()

	assert.NoError(err)

	relational.Replica = nil

	err = relational.Primary.Close()

	assert.NoError(err)

	relational.Primary = nil

	coreConfig.DatabaseConfiguration.Primary.Schema = baseSchema
	coreConfig.DatabaseConfiguration.Replica.Schema = baseSchema

	err = relational.Init(coreConfig.DatabaseConfiguration)

	assert.NoError(err)

	_, err = relational.Primary.Exec("DROP DATABASE " + tempSchema)

	if err != nil {
		logging.Warn(err)
	}

}

//freshStart creates a temporary database for the test
func freshStart(t *testing.T, coreConfig *config.CoreConfiguration) {

	tempSchema := coreConfig.DatabaseConfiguration.Primary.Schema + "_" + strings.ToLower(ids.NewSecureID())

	assert := assert.New(t)

	_, err := relational.Primary.Exec("CREATE DATABASE " + tempSchema)

	assert.NoError(err)

	coreConfig.DatabaseConfiguration.Primary.Schema = tempSchema
	coreConfig.DatabaseConfiguration.Replica.Schema = tempSchema

}

func testMigrations(t *testing.T, coreConfig *config.CoreConfiguration, loader loaders.ResourceLoader, schemaFiles []string) {

	assert := assert.New(t)

	err := relational.PreMigrate(loader, coreConfig.DatabaseConfiguration, schemaFiles)

	assert.NoError(err)

	//add post migrate call here

	changes, err := relational.Compare(loader, coreConfig.DatabaseConfiguration, schemaFiles)

	assert.NoError(err)

	assert.Equal(0, len(changes), "post migration database not in sync")

	if len(changes) > 0 {
		logging.LogJSONWithName("Changes:", changes)
	}

}
