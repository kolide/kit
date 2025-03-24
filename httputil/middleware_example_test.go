package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleChain() {
	h := Chain(
		annotate("one"),
		annotate("two"),
		annotate("three"),
	)(myHandler())

	srv := httptest.NewServer(h)
	defer srv.Close()

	resp, err := http.Get(srv.URL) // nolint:noctx
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Output:
	// annotate:  one
	// annotate:  two
	// annotate:  three
}

func annotate(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("annotate: ", s)
			next.ServeHTTP(w, r)
		})
	}
}

func myHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
