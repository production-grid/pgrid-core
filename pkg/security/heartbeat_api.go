package security

import (
	"net/http"
	"time"

	"github.com/production-grid/pgrid-core/pkg/applications"
	"github.com/production-grid/pgrid-core/pkg/httputils"
)

//HeartBeatResponse models a simple heartbeat response
type HeartBeatResponse struct {
	httputils.Acknowledgement
	CurrentTime string `json:"currentTime"`
	TenantType  string `json:"tenantType"`
}

//GetHeartbeat returns the current time
func GetHeartbeat(session applications.Session, w http.ResponseWriter, req *http.Request) {

	response := HeartBeatResponse{}
	response.Success = true
	response.CurrentTime = time.Now().String()
	response.TenantType = applications.CurrentApplication.TenantLingo.TenantPluralDefault

	httputils.SendJSON(response, w)

}
