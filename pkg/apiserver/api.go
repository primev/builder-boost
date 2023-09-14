package apiserver

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

const (
	defaultNamespace = "primev"
)

type searcherKey struct{}

type Service struct {
	*http.Server

	metricsRegistry *prometheus.Registry
	router          *http.ServeMux
	logger          *slog.Logger
}

func New() *Service {
	return &Service{}
}

func (a *Service) registerDebugEndpoints() {
	// register metrics handler
	a.router.Handle("/metrics", promhttp.HandlerFor(a.metricsRegistry, promhttp.HandlerOpts{}))

	// register pprof handlers
	a.router.Handle(
		"/debug/pprof",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Path += "/"
			http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
		}),
	)
	a.router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	a.router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	a.router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	a.router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	a.router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	a.router.Handle("/debug/pprof/{profile}", http.HandlerFunc(pprof.Index))
	a.router.Handle("/debug/vars", expvar.Handler())
}

func newMetrics(version string) (r *prometheus.Registry) {
	r = prometheus.NewRegistry()

	// register standard metrics
	r.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			Namespace: defaultNamespace,
		}),
		collectors.NewGoCollector(),
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Name:      "info",
			Help:      "builder-boost information.",
			ConstLabels: prometheus.Labels{
				"version": version,
			},
		}),
	)

	return r
}

func (a *Service) ChainHandlers(
	path string,
	handler http.Handler,
	mws ...func(http.Handler) http.Handler,
) {
	h := handler
	for i := len(mws) - 1; i > 0; i-- {
		h = mws[i](h)
	}
	a.router.Handle(path, h)
}

func (a *Service) RegisterMetricsCollectors(cs ...prometheus.Collector) {
	a.metricsRegistry.MustRegister(cs...)
}

