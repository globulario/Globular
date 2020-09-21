package monitoring_store

import (
	"context"
	"time"
)

/**
 * Represent a data store interface.
 */
type Store interface {
	// Alerts returns a list of all active alerts.
	Alerts(ctx context.Context) (string, error)
	// AlertManagers returns an overview of the current state of the Prometheus alert manager discovery.
	AlertManagers(ctx context.Context) (string, error)
	// CleanTombstones removes the deleted data from disk and cleans up the existing tombstones.
	CleanTombstones(ctx context.Context) error
	// Config returns the current Prometheus configuration.
	Config(ctx context.Context) (string, error)
	// DeleteSeries deletes data for a selection of series in a time range.
	DeleteSeries(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) error
	// Flags returns the flag values that Prometheus was launched with.
	Flags(ctx context.Context) (string, error)
	// LabelNames returns all the unique label names present in the block in sorted order.
	LabelNames(ctx context.Context) ([]string, string, error)
	// LabelValues performs a query for the values of the given label.
	LabelValues(ctx context.Context, label string, startTime int64, endTime int64) (string, string, error)
	// Query performs a query for the given time.
	Query(ctx context.Context, query string, ts time.Time) (string, string, error)
	// QueryRange performs a query for the given range.
	QueryRange(ctx context.Context, query string, startTime time.Time, endTime time.Time, step float64) (string, string, error)
	// Series finds series by label matchers.
	Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) (string, string, error)
	// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
	// under the TSDB's data directory and returns the directory as response.
	Snapshot(ctx context.Context, skipHead bool) (string, error)
	// Rules returns a list of alerting and recording rules that are currently loaded.
	Rules(ctx context.Context) (string, error)
	// Targets returns an overview of the current state of the Prometheus target discovery.
	Targets(ctx context.Context) (string, error)
	// TargetsMetadata returns metadata about metrics currently scraped by the target.
	TargetsMetadata(ctx context.Context, matchTarget string, metric string, limit string) (string, error)
}
