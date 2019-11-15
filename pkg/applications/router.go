package applications

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/production-grid/pgrid-core/pkg/httputils"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// Global server timeouts
const (
	// ReadTimouet is the time after which a connection is accepted until the request body is
	// fully read.
	ReadTimeout = 30 * time.Second

	// WriteTimoue is the time from the end of the request header to the end of
	// the response write.
	WriteTimeout = 300 * time.Second

	// IdleTimeout is a server side enforcement of keepalive timeout.
	IdleTimeout = 300 * time.Second
)

//Constants for crud access
const (
	CrudAccessRead   = "read"
	CrudAccessUpdate = "update"
	CrudAccessDelete = "delete"
)

// Route models route attributes common to all route types
type Route struct {
	Permission   string   //permission key
	AllRequired  []string //permission keys for multiple perms
	AnyRequired  []string
	Path         string
	TenantScoped bool //a tenant context is required for this route
	HandlerFunc  SecureHandlerFunc
	CORS         bool
}

// APIRoute models API specific route stuff
// API Routes paths must be relative to /api/<modulename
type APIRoute struct {
	Route
	Method string //http method
}

// ContentRoute models content route stuff
type ContentRoute struct {
	Route
}

// PermissionGroup models the three different ways of setting permission
type PermissionGroup struct {
	Permission    string   //permission key
	AllRequired   []string //permission keys for multiple perms
	AnyRequired   []string
	AuthorizeFunc AuthorizeFunc
}

// PagedResults models results with paging and sorting metadata
type PagedResults struct {
	Results
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
	PageSize  int `json:"pageSize"`
}

// Results models results without paging and sorting metadata
type Results struct {
	TotalResultCount int           `json:"totalResultCount"`
	VisibleResults   []interface{} `json:"visibleResults"`
}

// CrudResourcePermissions encapsulates permissions settings for crud resources
type CrudResourcePermissions struct {
	TenantScoped      bool
	ReadPermissions   PermissionGroup
	UpdatePermissions PermissionGroup
	DeletePermissions PermissionGroup
}

// AuthorizeFunc specifies an alternate function that can be used to enforce more granular data access
type AuthorizeFunc func(session *Session, req *http.Request, object interface{}) bool

//crudResourceRef is used internally to store router generated meta data
// about a resource
type crudResourceRef struct {
	Path        string
	Permissions CrudResourcePermissions
	Resource    CrudResource
}

//MetaDataPath returns the meta data path for the resource
func (ref *crudResourceRef) MetaDataPath() string {
	return ref.Path + "/md"
}

//MetaDataPath returns the meta data path for the resource
func (ref *crudResourceRef) OnePath() string {
	return ref.Path + "/{id}"
}

//MetaDataPath returns the meta data path for the resource
func (ref *crudResourceRef) DeletePath() string {
	return ref.Path + "/{id}"
}

//AllPath returns the path for returning all data elements for the resource
func (ref *crudResourceRef) AllPath() string {
	return ref.Path
}

//AllPath returns the path for returning all data elements for the resource
func (ref *crudResourceRef) UpdatePath() string {
	return ref.Path
}

// CrudResource models the contract for level 2 rest API resources
type CrudResource interface {
	Path() string
	Permissions() CrudResourcePermissions
	ToDTO(session *Session, req *http.Request, from interface{}) (interface{}, error)
	FromDTO(session *Session, req *http.Request, from interface{}, to interface{}) (interface{}, error)
	NewDTO(session *Session, req *http.Request) interface{}
	NewDomain(session *Session, req *http.Request) interface{}
	MetaData(session *Session, req *http.Request) CrudResourceMetaData
	All(session *Session, req *http.Request) ([]interface{}, error)
	One(session *Session, req *http.Request, id string) (interface{}, error)
	Update(session *Session, req *http.Request, dto interface{}, domain interface{}) (interface{}, error)
	Delete(session *Session, req *http.Request, domain interface{}) (bool, error)
}

// PagableCrudResource adds search, paging, and sort functionality to a basic API resource
type PagableCrudResource interface {
	CrudResource
	Page(session *Session, w http.ResponseWriter, req *http.Request)
	Sort(session *Session, w http.ResponseWriter, req *http.Request)
	Search(session *Session, w http.ResponseWriter, req *http.Request)
}

// DTO specifies the contract for DTO's
type DTO interface {
	Identifier() string
}

//SecureHandlerFunc is
type SecureHandlerFunc func(session Session, w http.ResponseWriter, req *http.Request)

// builds and starts the router
func initRouter(app *Application) error {

	logging.Debugln("Initializing HTTP Router...")

	app.Router = mux.NewRouter().StrictSlash(true)

	apiRouter := app.Router.PathPrefix("/api").Subrouter()

	for _, rc := range app.crudResources {
		logging.Debugf("Adding Resource Metadata Route %v: /api%v", http.MethodGet, rc.MetaDataPath())
		apiRouter.HandleFunc(rc.MetaDataPath(), metadataFunctionFor(rc)).Methods(http.MethodGet)
		logging.Debugf("Adding Resource All Route %v: /api%v", http.MethodGet, rc.AllPath())
		apiRouter.HandleFunc(rc.AllPath(), allFunctionFor(rc)).Methods(http.MethodGet)
		logging.Debugf("Adding Resource One Route %v: /api%v", http.MethodGet, rc.OnePath())
		apiRouter.HandleFunc(rc.OnePath(), oneFunctionFor(rc)).Methods(http.MethodGet)
		logging.Debugf("Adding Resource Update Route %v: /api%v", http.MethodPost, rc.UpdatePath())
		apiRouter.HandleFunc(rc.UpdatePath(), updateFunctionFor(rc)).Methods(http.MethodPost)
		logging.Debugf("Adding Resource Delete Route %v: /api%v", http.MethodDelete, rc.DeletePath())
		apiRouter.HandleFunc(rc.DeletePath(), deleteFunctionFor(rc)).Methods(http.MethodDelete)
	}

	for _, route := range app.apiRoutes {
		logging.Debugf("Adding API Route %v: /api%v", route.Method, route.Path)
		apiRouter.HandleFunc(route.Path, handlerFunctionFor(route.Route)).Methods(route.Method)
		if route.CORS {
			logging.Debugf("Adding API Route OPTIONS: /api%v", route.Path)
			apiRouter.HandleFunc(route.Path, SendCorsHeaders).Methods("OPTIONS")
		}
	}

	server := http.Server{
		Addr:         ":" + strconv.Itoa(app.CoreConfiguration.PortNumber),
		Handler:      Handler(app.Router),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	app.Server = &server

	return nil

}

func updateFunctionFor(rc crudResourceRef) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessUpdate, nil) {
			httputils.Send403(w)
			return
		}

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		rawDto := rc.Resource.NewDTO(&session, req)

		if httputils.ConsumeRequestBody(&rawDto, w, req) {

			dto := rawDto.(DTO)

			if err != nil {
				httputils.SendError(err, w)
				return
			}

			var domain interface{}
			var err error

			if dto.Identifier() != "" {
				domain, err = rc.Resource.One(&session, req, dto.Identifier())

				if err != nil {
					httputils.SendError(err, w)
					return
				}
				if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessUpdate, domain) {
					httputils.Send403(w)
					return
				}
			} else {
				domain = rc.Resource.NewDomain(&session, req)
			}

			updatedDomain, err := rc.Resource.FromDTO(&session, req, dto, domain)
			if err != nil {
				httputils.SendError(err, w)
				return
			}

			updatedDomain, err = rc.Resource.Update(&session, req, rawDto, updatedDomain)

			if err != nil {
				httputils.SendError(err, w)
				return
			}

			updatedDto, err := rc.Resource.ToDTO(&session, req, updatedDomain)
			if err != nil {
				httputils.SendError(err, w)
				return
			}

			httputils.SendJSON(updatedDto, w)

		}

	}

}

func deleteFunctionFor(rc crudResourceRef) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, nil) {
			httputils.Send403(w)
			return
		}

		vars := mux.Vars(req)

		domain, err := rc.Resource.One(&session, req, vars["id"])
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, domain) {
			httputils.Send403(w)
			return
		}

		ok, err := rc.Resource.Delete(&session, req, domain)
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !ok {
			httputils.SendJSON(httputils.Acknowledgement{Success: false, Error: "Delete Failed"}, w)
		} else {
			httputils.SendJSON(httputils.Acknowledgement{Success: true}, w)
		}

	}

}

func oneFunctionFor(rc crudResourceRef) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, nil) {
			httputils.Send403(w)
			return
		}

		vars := mux.Vars(req)

		domain, err := rc.Resource.One(&session, req, vars["id"])
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, domain) {
			httputils.Send403(w)
			return
		}
		dto, err := rc.Resource.ToDTO(&session, req, domain)
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		httputils.SendJSON(dto, w)

	}

}

func allFunctionFor(rc crudResourceRef) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, nil) {
			httputils.Send403(w)
			return
		}

		domainResults, err := rc.Resource.All(&session, req)
		if err != nil {
			httputils.SendError(err, w)
			return
		}

		visibleResults := make([]interface{}, 0)

		for _, domain := range domainResults {
			/*
				Seems redundant, but gives the AuthFunc the ability to block access
				to individual data elements.
			*/
			if isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, domain) {
				dto, err := rc.Resource.ToDTO(&session, req, domain)
				if err != nil {
					httputils.SendError(err, w)
					return
				}
				visibleResults = append(visibleResults, dto)
			}

		}

		results := PagedResults{
			Page:      1,
			PageCount: 1,
			PageSize:  len(visibleResults),
		}

		results.TotalResultCount = len(visibleResults)
		results.VisibleResults = visibleResults

		httputils.SendJSON(results, w)

	}

}

func metadataFunctionFor(rc crudResourceRef) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isCrudAccessAuthorized(&rc, req, &session, CrudAccessRead, nil) {
			httputils.Send403(w)
			return
		}

		md := rc.Resource.MetaData(&session, req)

		httputils.SendJSON(md, w)

	}

}

func handlerFunctionFor(route Route) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		session, err := resolveSession(req)

		if err != nil {
			httputils.SendError(err, w)
			return
		}

		if !isAuthorized(&route, req, &session) {
			httputils.Send403(w)
			return
		}

		route.HandlerFunc(session, w, req)

	}

}

func isCrudAccessAuthorized(rc *crudResourceRef, req *http.Request, session *Session, accessType string, domain interface{}) bool {

	var permissionGroup PermissionGroup

	switch accessType {
	case CrudAccessDelete:
		permissionGroup = rc.Permissions.DeletePermissions
	case CrudAccessUpdate:
		permissionGroup = rc.Permissions.UpdatePermissions
	default:
		permissionGroup = rc.Permissions.ReadPermissions
	}

	if rc.Permissions.TenantScoped && session.TenantID == "" {
		logging.Warnf("resource %v requires tenant context", rc.Path)
		return false
	}

	if permissionGroup.Permission != "" {
		if !session.HasPermission(permissionGroup.Permission) {
			return false
		}
	}

	if (permissionGroup.AllRequired != nil) && (len(permissionGroup.AllRequired) > 0) {
		for _, permCode := range permissionGroup.AllRequired {
			if !session.HasPermission(permCode) {
				return false
			}
		}
	}

	if (permissionGroup.AnyRequired != nil) && (len(permissionGroup.AnyRequired) > 0) {
		for _, permCode := range permissionGroup.AnyRequired {
			if session.HasPermission(permCode) {
				return isAuthorizeFuncApproved(&permissionGroup, session, req, domain)
			}
		}
		return false
	}

	return isAuthorizeFuncApproved(&permissionGroup, session, req, domain)

}

func isAuthorizeFuncApproved(permGroup *PermissionGroup, session *Session, req *http.Request, object interface{}) bool {

	if permGroup.AuthorizeFunc == nil {
		return true
	}

	return permGroup.AuthorizeFunc(session, req, object)

}

func isAuthorized(route *Route, req *http.Request, session *Session) bool {

	if route.TenantScoped && session.TenantID == "" {
		logging.Warnf("route %v requires tenant context", route.Path)
		return false
	}

	if route.Permission != "" {
		if !session.HasPermission(route.Permission) {
			return false
		}
	}

	if (route.AllRequired != nil) && (len(route.AllRequired) > 0) {
		for _, permCode := range route.AllRequired {
			if !session.HasPermission(permCode) {
				return false
			}
		}
	}

	if (route.AnyRequired != nil) && (len(route.AnyRequired) > 0) {
		for _, permCode := range route.AnyRequired {
			if session.HasPermission(permCode) {
				return true
			}
		}
		return false
	}

	return true

}

// Handler mutates server responses. It is used to add standard headers or wrap
// the default handler with additional layers of functionality, including
// request logging.
func Handler(next http.Handler) http.Handler {
	stsHeader := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		next.ServeHTTP(res, req)
	})

	return handlers.ProxyHeaders(
		handlers.CombinedLoggingHandler(
			os.Stdout,
			stsHeader,
		),
	)
}

// SendCorsHeaders sets the property headers for cross origin support.
func SendCorsHeaders(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Timestamp, Nonce, Authorization, Content-Type")

}
