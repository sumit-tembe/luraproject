package gorilla

import (
	"net/http"
	"strings"

	gorilla "github.com/gorilla/mux"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/router/mux"
	"github.com/luraproject/lura/v2/transport/http/server"

	customMux "github.com/sumit-tembe/luraproject/mux"
	customProxy "github.com/sumit-tembe/luraproject/proxy"
)

// DefaultConfig returns the struct that collects the parts the router should be builded from
func DefaultConfig(pf proxy.Factory, logger logging.Logger) mux.Config {
	proxy.NewRequestBuilderMiddleware = customProxy.NewRequestBuilderMiddleware
	return mux.Config{
		Engine:         gorillaEngine{gorilla.NewRouter()},
		Middlewares:    []mux.HandlerMiddleware{},
		HandlerFactory: mux.CustomEndpointHandler(customMux.NewRequestBuilder(gorillaParamsExtractor)),
		ProxyFactory:   pf,
		Logger:         logger,
		DebugPattern:   "/__debug/{params}",
		RunServer:      server.RunServer,
	}
}

func gorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	for key, value := range gorilla.Vars(r) {
		params[strings.Title(key)] = value
	}
	return params
}

type gorillaEngine struct {
	r *gorilla.Router
}

// Handle implements the mux.Engine interface from the lura router package
func (g gorillaEngine) Handle(pattern, method string, handler http.Handler) {
	if strings.HasSuffix(pattern, "*") {
		path := strings.TrimSuffix(pattern, "*")
		g.r.PathPrefix(path).Handler(handler).Name(pattern)
		return
	}
	g.r.Handle(pattern, handler).Methods(method).Name(pattern)
}

// ServeHTTP implements the http:Handler interface from the stdlib
func (g gorillaEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.r.ServeHTTP(mux.NewHTTPErrorInterceptor(w), r)
}
