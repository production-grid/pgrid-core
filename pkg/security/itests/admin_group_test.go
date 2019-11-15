package itests

import (
	"testing"

	"github.com/production-grid/pgrid-core/pkg/security"
	"github.com/production-grid/pgrid-core/pkg/testutil"
)

func TestAdminGroupLifecycle(t *testing.T) {

	testutil.TestDomainLifeCycle(t, &security.AdminGroup{}, &security.AdminGroupFinder{})

}
