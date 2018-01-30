package debug

import (
	"context"
	"net"
	"net/http"
	nhpprof "net/http/pprof"
	"net/url"
	"runtime/pprof"
	"strings"

	"github.com/alecthomas/template"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

// Server is the debug server struct. It should be created through StartServer.
type Server struct {
	serv      *http.Server
	addr      string
	authToken string
	logger    log.Logger
	prefix    string
}

// Option is the functional option type for Server.
type Option func(*Server)

// WithAddr sets the address to bind to.
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithAuthToken sets the auth token to use. If it is unset, there is no auth.
func WithAuthToken(token string) Option {
	return func(s *Server) {
		s.authToken = token
	}
}

// WithLogger sets the logger to use.
func WithLogger(logger log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithPrefix sets the URL prefix to use.
func WithPrefix(prefix string) Option {
	return func(s *Server) {
		s.prefix = prefix
	}
}

// StartServer creates and starts a new debug server using the provided
// functional Options.
func StartServer(opts ...Option) (*Server, error) {
	s := &Server{
		addr:      ":63809",
		authToken: "",
		logger:    log.NewNopLogger(),
		prefix:    "/",
	}
	for _, opt := range opts {
		opt(s)
	}

	m := http.NewServeMux()
	h := handler(s.authToken, s.logger)
	if s.authToken != "" {
		h = authHandler(s.authToken, s.logger)
	}
	m.Handle(s.prefix, http.StripPrefix(s.prefix, h))
	s.serv = &http.Server{
		Handler: m,
	}

	return s, s.Start()
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errors.Wrap(err, "opening socket")
	}

	go func() {
		defer l.Close()
		if err := s.serv.Serve(l); err != nil && err != http.ErrServerClosed {
			level.Info(s.logger).Log("msg", "debug server failed", "err", err)
		}
	}()

	url := url.URL{
		Scheme:   "http",
		Host:     l.Addr().String(),
		Path:     s.prefix,
		RawQuery: "token=" + s.authToken,
	}
	addr := url.String()
	level.Info(s.logger).Log(
		"msg", "debug server started",
		"addr", addr,
	)

	return nil
}

// Shutdown stops the running debug server.
func (s *Server) Shutdown() error {
	err := s.serv.Shutdown(context.Background())
	return errors.Wrap(err, "shutting down server")
}

// The below handler code is adapted from MIT licensed github.com/e-dard/netbug
func handler(token string, logger log.Logger) http.HandlerFunc {
	info := struct {
		Profiles []*pprof.Profile
		Token    string
	}{
		Profiles: pprof.Profiles(),
		Token:    url.QueryEscape(token),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		switch name {
		case "":
			// Index page.
			if err := indexTmpl.Execute(w, info); err != nil {
				level.Info(logger).Log(
					"msg", "error rendering debug template",
					"err", err,
				)
				return
			}
		case "cmdline":
			nhpprof.Cmdline(w, r)
		case "profile":
			nhpprof.Profile(w, r)
		case "trace":
			nhpprof.Trace(w, r)
		case "symbol":
			nhpprof.Symbol(w, r)
		default:
			// Provides access to all profiles under runtime/pprof
			nhpprof.Handler(name).ServeHTTP(w, r)
		}
	}
}

// authHandler wraps the basic handler, checking the auth token.
func authHandler(token string, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("token") == token {
			handler(token, logger).ServeHTTP(w, r)
		} else {
			http.Error(w, "Request must include valid token.", http.StatusUnauthorized)
		}
	}
}

var indexTmpl = template.Must(template.New("index").Parse(`<html>
  <head>
    <title>Debug Information</title>
  </head>
  <br>
  <body>
    Profiles:<br>
    <table>
    {{range .Profiles}}
      <tr><td align=right>{{.Count}}<td><a href="{{.Name}}?debug=1&token={{$.Token}}">{{.Name}}</a>
    {{end}}
    <tr><td align=right><td><a href="profile?token={{.Token}}">CPU</a>
    <tr><td align=right><td><a href="trace?seconds=5&token={{.Token}}">5-second trace</a>
    <tr><td align=right><td><a href="trace?seconds=30&token={{.Token}}">30-second trace</a>
    </table>
    <br>
    Debug information:<br>
    <table>
      <tr><td align=right><td><a href="cmdline?token={{.Token}}">cmdline</a>
      <tr><td align=right><td><a href="symbol?token={{.Token}}">symbol</a>
    <tr><td align=right><td><a href="goroutine?debug=2&token={{.Token}}">full goroutine stack dump</a><br>
    <table>
  </body>
</html>`))
