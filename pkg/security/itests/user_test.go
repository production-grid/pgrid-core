package itests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/production-grid/pgrid-core/pkg/database/relational"
	"github.com/production-grid/pgrid-core/pkg/logging"
	"github.com/production-grid/pgrid-core/pkg/security"
	"github.com/production-grid/pgrid-core/pkg/testutil"
)

func TestUserLifecycle(t *testing.T) {

	assert := assert.New(t)

	testutil.StartTestApplication(t)

	user := security.User{}
	testutil.PopulateTestData(&user)

	id, err := user.Save()

	assert.NoError(err)
	assert.NotEmpty(id)

	userFinder := security.UserFinder{}

	savedUser, err := userFinder.FindByID(relational.REPLICA, id)

	assert.NoError(err)
	assert.NotNil(savedUser)

	logging.LogJSON(savedUser)

	testutil.AssertEquivalent(t, *savedUser, user)

}
