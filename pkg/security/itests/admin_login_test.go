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

	//invalid user attempt
	ack := httputils.Acknowledgement{}

	loginRequest := security.LoginRequest{
		EmailAddress: "isnogood@productiongrid.com",
		Password:     "test123",
	}

	testutil.PostJSON(t, "/api/security/login", loginRequest, &ack)

	assert.False(ack.Success)

	//bad password
	loginRequest = security.LoginRequest{
		EmailAddress: "devops@productiongrid.com",
		Password:     "isnogood",
	}

	ack = httputils.Acknowledgement{}

	testutil.PostJSON(t, "/api/security/login", loginRequest, &ack)

	assert.False(ack.Success)

	//good login attempt
	loginRequest = security.LoginRequest{
		EmailAddress: "devops@productiongrid.com",
		Password:     "test123",
	}

	ack = httputils.Acknowledgement{}

	testutil.PostJSON(t, "/api/security/login", loginRequest, &ack)

	assert.True(ack.Success)

	session := security.SessionDTO{}

	testutil.GetJSONWithSessionKey(t, "/api/security/session", ack.ID, &session)

	assert.NotEmpty(session.UserID)

	//logout
	ack = httputils.Acknowledgement{}

	testutil.PostJSON(t, "/api/security/logout", nil, &ack)

	assert.True(ack.Success)

	//ensure the session has been destroyed
	session = security.SessionDTO{}

	testutil.GetJSONWithSessionKey(t, "/api/security/session", ack.ID, &session)

	assert.Empty(session.UserID)

}
