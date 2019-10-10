package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/production-grid/pgrid-core/pkg/database/relational"
)

// TestDomainLifeCycle encapsulates generic entity stack testing
func TestDomainLifeCycle(t *testing.T, entity relational.Entity, finder relational.EntityFinder) {

	assert := assert.New(t)

	StartTestApplication(t)

	PopulateTestData(entity)

	id, err := entity.Save()

	assert.NoError(err)
	assert.NotEmpty(id)

	savedEntity, err := finder.FindInterfaceByID(relational.REPLICA, id)

	assert.NoError(err)
	assert.NotNil(savedEntity)

	AssertEquivalent(t, savedEntity, entity)

	PopulateTestData(entity)

	id, err = entity.Save()

	assert.NoError(err)
	assert.NotEmpty(id)

	savedEntity, err = finder.FindInterfaceByID(relational.REPLICA, id)

	assert.NoError(err)
	assert.NotNil(savedEntity)

	AssertEquivalent(t, savedEntity, entity)

	err = entity.Delete()
	assert.NoError(err)

	savedEntity, err = finder.FindInterfaceByID(relational.REPLICA, id)

	assert.Error(err)
	assert.Nil(savedEntity)

}
