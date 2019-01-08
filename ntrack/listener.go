package ntrack

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type trackingListener struct {
	net.Listener
	stats *Stats
}

func NewInstrumentedListener(lis net.Listener) (net.Listener, *Stats) {
	listenerStats := &Stats{}
	listenerStats.init()

	return &trackingListener{
		Listener: lis,
		stats:    listenerStats,
	}, listenerStats
}

func (tl *trackingListener) Accept() (net.Conn, error) {
	conn, err := tl.Listener.Accept()
	stats.RecordWithTags(context.TODO(), []tag.Mutator{tag.Upsert(tl.stats.TagSuccess, fmt.Sprintf("%v", err == nil))}, tl.stats.ListenerAccepted.M(1))
	if err != nil {
		return nil, errors.Wrap(err, "accept from base listener")
	}

	atomic.AddInt64(&tl.stats.openConnections, 1)
	open := atomic.LoadInt64(&tl.stats.openConnections)
	stats.Record(context.TODO(), tl.stats.OpenConnections.M(open))
	return &serverConn{Conn: conn, stats: tl.stats}, nil
}

type serverConn struct {
	net.Conn
	stats *Stats
}

func (sc *serverConn) Close() error {
	err := sc.Conn.Close()
	atomic.AddInt64(&sc.stats.openConnections, -1)
	open := atomic.LoadInt64(&sc.stats.openConnections)
	stats.Record(context.TODO(),
		sc.stats.OpenConnections.M(open),
		sc.stats.LifetimeClosedConnections.M(1),
	)
	return errors.Wrap(err, "close server conn")
}

type Stats struct {
	ListenerAccepted          *stats.Int64Measure
	LifetimeClosedConnections *stats.Int64Measure
	OpenConnections           *stats.Int64Measure
	openConnections           int64

	TagSuccess tag.Key

	views []*view.View
}

func (s *Stats) init() {
	s.ListenerAccepted = stats.Int64("ntrack/listener/accepts", "The number of Accept calls on the net.Listener", stats.UnitDimensionless)
	s.LifetimeClosedConnections = stats.Int64("ntrack/listener/closed", "The number of Close calls on the net.Listener", stats.UnitDimensionless)
	s.OpenConnections = stats.Int64("ntrack/listener/open", "The number of Open connections from the net.Listener", stats.UnitDimensionless)

	s.TagSuccess, _ = tag.NewKey("success")

	tags := []tag.Key{s.TagSuccess}
	s.views = append(s.views, viewFromStat(s.ListenerAccepted, tags, view.Count()))
	s.views = append(s.views, viewFromStat(s.OpenConnections, nil, view.LastValue()))
	s.views = append(s.views, viewFromStat(s.LifetimeClosedConnections, nil, view.Count()))
}

func viewFromStat(ss *stats.Int64Measure, tags []tag.Key, agg *view.Aggregation) *view.View {
	return &view.View{
		Name:        ss.Name(),
		Measure:     ss,
		Description: ss.Description(),
		TagKeys:     tags,
		Aggregation: agg,
	}
}
