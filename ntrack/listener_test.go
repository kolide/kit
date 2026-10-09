package ntrack

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
)

func TestListener(t *testing.T) {
	var tests = []struct {
		viewName         string
		disableKeepalive bool
		expectedValue    int64
	}{
		{
			viewName:         "ntrack/listener/accepts",
			disableKeepalive: true,
			expectedValue:    5,
		},
		{
			viewName:         "ntrack/listener/accepts",
			disableKeepalive: false,
			expectedValue:    1,
		},
		{
			viewName:         "ntrack/listener/closed",
			disableKeepalive: true,
			expectedValue:    5,
		},
		{
			viewName:         "ntrack/listener/open",
			disableKeepalive: true,
			expectedValue:    0,
		},
		{
			viewName:         "ntrack/listener/open",
			disableKeepalive: false,
			expectedValue:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.viewName, func(t *testing.T) {
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			ilis, stats := NewInstrumentedListener(lis)
			registerViewByName(t, tt.viewName, stats, false)

			testClientConnections(t, ilis, tt.disableKeepalive)

			rows, err := view.RetrieveData(tt.viewName)
			require.NoError(t, err)

			switch data := rows[0].Data.(type) {
			case *view.CountData:
				assert.Equal(t, tt.expectedValue, data.Value)
			case *view.LastValueData:
				assert.Equal(t, float64(tt.expectedValue), data.Value)
			}
			registerViewByName(t, tt.viewName, stats, true)
		})
	}
}

func registerViewByName(t *testing.T, name string, stats *Stats, unregister bool) {
	var v *view.View
	switch name {
	case "ntrack/listener/accepts":
		v = stats.ListenerAcceptedView
	case "ntrack/listener/open":
		v = stats.OpenConnectionsView
	case "ntrack/listener/closed":
		v = stats.LifetimeClosedConnectionsView
	}
	if unregister {
		view.Unregister(v)
	} else {
		view.Register(v)
	}
}

func testClientConnections(t *testing.T, lis net.Listener, disableKeepalive bool) {
	t.Helper()

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Fatal(err)
		}
	}()

	tr := &http.Transport{DisableKeepAlives: disableKeepalive}
	client := &http.Client{Transport: tr}

	requestCount := 5
	for i := 0; i < requestCount; i++ {
		resp, err := client.Get(fmt.Sprintf("http://%s", lis.Addr()))
		require.NoError(t, err)
		resp.Body.Close()
	}

}
