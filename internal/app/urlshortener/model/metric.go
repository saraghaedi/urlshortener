package model

import (
	"time"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/metric"

	prom "github.com/saraghaedi/urlshortener/pkg/prometheus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelRepoName      = "repo_name"
	labelRepoMethod    = "repo_method"
	errorIncrementStep = 1
)

// Metrics keeps global Prometheus metrics.
type Metrics struct {
	ErrCounter *prometheus.CounterVec
	Histogram  *prometheus.HistogramVec
}

// nolint:gochecknoglobals
var (
	metrics = Metrics{
		ErrCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metric.Namespace,
				Name:      "repo_error_total",
				Help:      "Repo error total.",
			}, []string{labelRepoName, labelRepoMethod}),

		Histogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: metric.Namespace,
			Name:      "repo_duration_total",
			Help:      "Repo duration total.",
			Buckets:   prom.HistogramBuckets,
		}, []string{labelRepoName, labelRepoMethod}),
	}
)

func (m Metrics) report(repoName, methodName string, startTime time.Time, err error) {
	if err != nil {
		m.ErrCounter.With(prometheus.Labels{labelRepoName: repoName, labelRepoMethod: methodName}).Add(errorIncrementStep)

		return
	}

	m.Histogram.With(prometheus.Labels{labelRepoName: repoName, labelRepoMethod: methodName}).
		Observe(time.Since(startTime).Seconds())
}
