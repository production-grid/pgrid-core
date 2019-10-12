package itests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/production-grid/pgrid-core/pkg/httputils"
	"github.com/production-grid/pgrid-core/pkg/security"
	"github.com/production-grid/pgrid-core/pkg/testutil"
)

func TestAdminLogin(t *testing.T) {

	assert := assert.New(t)

	testutil.StartTestApplication(t)

	loginRequest := security.LoginRequest{
		EmailAddress: "devops@productiongrid.com",
		Password:     "test123",
	}

	ack := httputils.Acknowledgement{}

	testutil.PostJSON(t, "/api/security/login", loginRequest, &ack)

	assert.True(ack.Success)

}
