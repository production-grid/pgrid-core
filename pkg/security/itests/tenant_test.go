package itests

import (
	"testing"

	"github.com/production-grid/pgrid-core/pkg/security"
	"github.com/production-grid/pgrid-core/pkg/testutil"
)

func TestTenantLifecycle(t *testing.T) {

	testutil.TestDomainLifeCycle(t, &security.Tenant{}, &security.TenantFinder{})

}
