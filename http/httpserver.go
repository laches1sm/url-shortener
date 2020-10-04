package http

import (
	"github.com/laches1sm/url_shortener/adapters"
	"github.com/laches1sm/url_shortener/infrastructure"
	"log"
	"net/http"
)

const (
	ServerPort     = ":7000"
	ShortenEndpoint = "/shorten"
)

// ParrotServer is an interface to an HTTP server which handles requests for parrots
type ParrotServer struct {
	Mux               *http.ServeMux
	Logger            log.Logger
	UrlAdapter *adapters.UrlShortenerAdapter
	Infra *infrastructure.UrlInfra
}

// SetupRoutes configures the routes of the API
func (srv *ParrotServer) SetupRoutes() {
	srv.Mux.Handle(ShortenEndpoint, http.HandlerFunc(srv.UrlAdapter.ShortenURL))
	srv.Mux.Handle("/{url}", http.HandlerFunc(srv.UrlAdapter.RedirectUrl))
}

// Start sets up the HTTP webserver to listen and handle traffic. It
// takes the port number to listen on as a parameter in the form ":PORT_NUMBER"
func (srv *ParrotServer) Start(port string) error {
	return http.ListenAndServe(port, srv.Mux)
}

// NewParrotServer returns an instance of a configured ParrotServer
func NewParrotServer(logger log.Logger, adapter *adapters.UrlShortenerAdapter, infra *infrastructure.UrlInfra) *ParrotServer {
	httpServer := &ParrotServer{
		Mux:               http.NewServeMux(),
		Logger:            logger,
		UrlAdapter: adapter,
		Infra: infra,
	}
	return httpServer
}
