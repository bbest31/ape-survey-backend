package server

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/apesurvey/ape-survey-backend/v2/utils"
	"github.com/joho/godotenv"
)

var (
	//defaultWriteTimeout is the max time the client has to write a response
	defaultWriteTimeout = 15 * time.Second
	defaultReadTimeout  = 15 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

type HTTPServer struct {
	Host         string
	Port         string
	WriteTimeout time.Duration // max time to wait to write a response
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Router       http.Handler
	TLS          bool
	TLSConfig    *tls.Config
	CertFile     string
	KeyFile      string
}

type ServerOptions struct {
	WriteTimeout time.Duration // max time to wait to write a response
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	TLS          bool
	TLSConfig    tls.Config
	CertFile     string
	KeyFile      string
}

// InitEnvironment grabs the environment variables necessary to make the http server.
func (s *HTTPServer) initEnvironment() {
	err := godotenv.Load()
	if err != nil {
		s.Host = "localhost"
		s.Port = "8080"
		return
	}

	host, ok := os.LookupEnv("HOST")
	if !ok || utils.IsEmptyString(host) {
		s.Host = "localhost"
	} else {
		s.Host = host
	}

	port, ok := os.LookupEnv("PORT")
	if !ok || utils.IsEmptyString(port) {
		s.Port = "8080"
	} else {
		s.Port = port
	}
}

// DefaultServer creates a new HTTPServer with the default options.
func DefaultServer() (*HTTPServer, error) {
	return NewServer(ServerOptions{})

}

// NewServer create a new server with default settings.
func NewServer(opts ServerOptions) (*HTTPServer, error) {
	server := &HTTPServer{
		WriteTimeout: opts.WriteTimeout,
		ReadTimeout:  opts.ReadTimeout,
		IdleTimeout:  opts.IdleTimeout,
		TLS:          opts.TLS,
		TLSConfig:    &opts.TLSConfig,
		CertFile:     opts.CertFile,
		KeyFile:      opts.KeyFile,
		Router:       http.DefaultServeMux,
	}

	if opts.WriteTimeout == 0 {
		server.WriteTimeout = defaultWriteTimeout
	}

	if opts.ReadTimeout == 0 {
		server.ReadTimeout = defaultReadTimeout
	}

	if opts.IdleTimeout == 0 {
		server.IdleTimeout = defaultIdleTimeout
	}

	server.initEnvironment()

	return server, nil
}

// SetRouter sets the route handler for the client. If nil is passed we use the default serve mux.
func (s *HTTPServer) SetRouter(router http.Handler) {

	if router == nil {
		s.Router = http.DefaultServeMux
	} else {
		s.Router = router
	}
}

// ListenAndServe will setup and launch an http server based on the client options.
func (s *HTTPServer) ListenAndServe() {

	srv := &http.Server{
		WriteTimeout: s.WriteTimeout,
		ReadTimeout:  s.ReadTimeout,
		IdleTimeout:  s.IdleTimeout,
		Addr:         s.Host + ":" + s.Port,
		Handler:      s.Router,
		TLSConfig:    s.TLSConfig,
	}

	log.Printf("Server listening and running on %v:%v\n", s.Host, s.Port)
	if s.TLS {
		log.Fatalln(srv.ListenAndServeTLS(s.CertFile, s.KeyFile))
	} else {
		log.Fatalln(srv.ListenAndServe())
	}
}
