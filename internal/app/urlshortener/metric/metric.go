package metric

import (
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// Namespace is the Prometheus metric namespace variable.
	Namespace = "urlshortener"

	labelDbName = "db_name"
)

// Metrics keeps global Prometheus metrics.
type Metrics struct {
	DbConnectionStatus *prometheus.GaugeVec
}

// nolint:gochecknoglobals
var metrics = Metrics{
	DbConnectionStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "db_connection_status",
		Help:      "Database connection status.",
	}, []string{labelDbName}),
}

// ReportDbStatus reports status of database connection to the Prometheus.
func ReportDbStatus(db *gorm.DB, dbName string) {
	// 1 means query is ok and 0 means query is not ok
	status := 1
	if err := db.Exec("SELECT 1;").Error; err != nil {
		status = 0
	}

	metrics.DbConnectionStatus.With(prometheus.Labels{labelDbName: dbName}).Set(float64(status))
}
