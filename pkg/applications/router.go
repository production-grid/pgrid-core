package applications

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

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

// Route models route attributes common to all route types
type Route struct {
	Permission   string   //permission key
	AllRequired  []string //permission keys for multiple perms
	AnyRequired  []string
	Path         string
	TenantScoped bool //a tenant context is required for this route
	HandlerFunc  http.HandlerFunc
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

// builds and starts the router
func initRouter(app *Application) error {

	logging.Debugln("Initializing HTTP Router...")

	app.Router = mux.NewRouter().StrictSlash(true)

	apiRouter := app.Router.PathPrefix("/api").Subrouter()

	for _, route := range app.APIRoutes {
		logging.Debugf("Adding API Route %v: /api%v", route.Method, route.Path)
		apiRouter.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method)
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

// Handler mutates server responses. It is used to add standard headers or wrap
// the default handler with additional layers of functionality, including
// request logging.
func Handler(next http.Handler) http.Handler {
	stsHeader := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		next.ServeHTTP(res, req)
	})
	// External handlers should be added here

	// Switch back to gorilla fork once this issue is merged:
	// https://github.com/gorilla/handlers/issues/141
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
