// Example program using kit/httputil to launch a web server. Content
// is served via https, with a redirect on http. It uses autocert to
// fetch a letsencrypt cert (from staging).
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kolide/kit/httputil"
	"github.com/oklog/run"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	var (
		flHostName = flag.String("hostname", "", "External hostname for service. Needed for autocert")
		flCertDir  = flag.String("certdir", "", "Directory to store certificates in")
	)
	flag.Parse()

	if *flHostName == "" {
		fmt.Println("Must specify hostname")
		os.Exit(1)
	}

	if *flCertDir == "" {
		fmt.Println("Must specify certdir")
		os.Exit(1)
	}

	m, err := httputil.NewAutocertManager(
		autocert.DirCache(*flCertDir),
		[]string{*flHostName},
		httputil.WithLetsEncryptStaging(),
	)
	if err != nil {
		panic(err)
	}

	var g run.Group
	{
		srv := httputil.NewServer(
			":443",
			stringMiddleware("secure!", nil),
			httputil.WithTLSConfig(&tls.Config{GetCertificate: m.GetCertificate}),
		)

		g.Add(func() error {
			fmt.Println("Starting port 443")
			return srv.ListenAndServeTLS("", "")
		}, func(err error) {
			srv.Close()
			return
		})
	}
	{
		srv := httputil.NewServer(
			":80",
			m.HTTPHandler(httputil.RedirectToSecureHandler()),
			httputil.WithReadTimeout(5*time.Second),
			httputil.WithWriteTimeout(5*time.Second),
		)

		g.Add(func() error {
			fmt.Println("Starting port 80")
			return srv.ListenAndServe()
		}, func(err error) {
			srv.Close()
			return
		})
	}

	{
		// this actor handles an os interrupt signal and terminates the server.
		sig := make(chan os.Signal, 1)
		g.Add(func() error {
			signal.Notify(sig, os.Interrupt)
			<-sig
			fmt.Println("beginning shutdown")
			return nil
		}, func(err error) {
			fmt.Println("process interrupted")
			close(sig)
		})
	}

	if err := g.Run(); err != nil {
		panic(err)
	}
}

func stringMiddleware(s string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, s)
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
