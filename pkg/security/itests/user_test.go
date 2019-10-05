package itests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/production-grid/pgrid-core/pkg/security"
	"github.com/production-grid/pgrid-core/pkg/testutil"
)

func TestUserLifecycle(t *testing.T) {

	assert := assert.New(t)

	testutil.StartTestApplication(t)

	user := security.User{}
	user.EMail = testutil.GenerateTestEmail(12)

	id, err := user.Save()

	assert.NoError(err)
	assert.NotEmpty(id)

}
