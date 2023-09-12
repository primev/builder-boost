package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
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

type API struct {
	*http.Server

	metricsRegistry *prometheus.Registry
	router          *http.ServeMux
	logger          *slog.Logger
}

func NewAPI() *API {
	return &API{}
}

func (a *API) registerDebugEndpoints() {
	// register metrics handler
	a.router.Handle("/metrics", promhttp.HandlerFor(a.MetricsRegistry(), promhttp.HandlerOpts{}))

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

func (a *API) MetricsRegistry() *prometheus.Registry {
	return a.metricsRegistry
}

func (a *API) Router() *http.ServeMux {
	return a.router
}

type statusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteResponse(w http.ResponseWriter, code int, message any) error {
	var b bytes.Buffer
	switch message.(type) {
	case string:
		err := json.NewEncoder(&b).Encode(statusResponse{Code: code, Message: message.(string)})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to encode status response: %w", err)
		}
	default:
		err := json.NewEncoder(&b).Encode(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to encode response: %w", err)
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, b.String())
	return nil
}

func MethodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			WriteResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handler(w, r)
	}
}

func BindJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var body T

	if r.Body == nil {
		return body, errors.New("no body")
	}
	defer r.Body.Close()

	return body, json.NewDecoder(r.Body).Decode(&body)
}
