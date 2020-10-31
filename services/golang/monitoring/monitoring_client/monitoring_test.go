package monitoring_client

import (
	"fmt"
	"log"
	"testing"

	//	"time"

	"github.com/globulario/Globular/monitoring/monitoring_client"
)

// Set the correct addresse here as needed.
var (
	client = monitoring_client.NewMonitoring_Client("localhost", "monitoring_server")
)

// First test create a fresh new connection...
func TestMonitoring(t *testing.T) {
	fmt.Println("Monitoring test.")
}

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")
	// err := client.CreateConnection("test", "127.0.0.1", 0, 9090)
	err := client.CreateConnection("test", "steve_pc", 0, 9090)
	if err != nil {
		log.Println("Fail to create a new connection", err)
		t.Fail()
	}
}

// Test getting the configurations infromations.
func TestGetConfig(t *testing.T) {
	fmt.Println("Get configuration test.")
	config, err := client.Config("test")
	if err != nil {
		log.Println("Fail to get test config", err)
		t.Fail()
	}

	log.Println(config)
}

/*
// Test Alerts.
func TestAlerts(t *testing.T) {
	fmt.Println("Get a alerts")
	alerts, err := client.Alerts("test")
	if err != nil {
		log.Println("fail to get alerts", err)
		t.Fail()
	}

	log.Println(alerts)
}

// Test AlertManagers.
func TestAlertManagers(t *testing.T) {
	fmt.Println("Get a alerts managers")
	alerts, err := client.AlertManagers("test")
	if err != nil {
		log.Println("fail to get alerts manager", err)
		t.Fail()
	}

	log.Println(alerts)
}

func TestCleanTombstones(t *testing.T) {
	fmt.Println("Clean Tombstones")
	err := client.CleanTombstones("test")
	if err != nil {
		log.Println("fail to Clean Tombstones", err)
		t.Fail()
	}
}

func TestDeleteSeries(t *testing.T) {

}

func TestDeleteFlags(t *testing.T) {
	fmt.Println("Get Flags")
	flags, err := client.Flags("test")
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}
	log.Println(flags)
}

// Test get labels names.
func TestLabelNames(t *testing.T) {
	fmt.Println("Get label names")
	names, warnings, err := client.LabelNames("test")
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}
	log.Println(names)
	log.Println(warnings)
}

// Get a label value.
func TestLabelValues(t *testing.T) {
	fmt.Println("Get label values")
	values, warnings, err := client.LabelValues("test", "__name__")
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}
	log.Println(values)
	log.Println(warnings)
}

func TestQuery(t *testing.T) {
	fmt.Println("Run a query")
	ts := time.Now().Unix() - 1000
	// q := "rate(prometheus_tsdb_head_chunks_created_total[1m])"
	q := "histogram_quantile(0.95, sum(rate(prometheus_http_request_duration_seconds_bucket[5m])) by (le))"
	values, warnings, err := client.Query("test", q, float64(ts))
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}
	log.Println(values)
	log.Println(warnings)
}

func TestQueryRange(t *testing.T) {
	fmt.Println("Query a range of values")

	startTime := time.Now().Unix() - 1000
	endTime := time.Now().Unix() - 100
	step := float64(10000) // 10s

	q := "prometheus_target_interval_length_seconds"

	values, warnings, err := client.QueryRange("test", q, float64(startTime), float64(endTime), step)
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}
	log.Println(values)
	log.Println(warnings)
}

func TestSeries(t *testing.T) {
	fmt.Println("Get Series by labels")

	startTime := time.Now().Unix() - 10000
	endTime := time.Now().Unix()

	values, warnings, err := client.Series("test", []string{"prometheus_target_interval_length_seconds"}, float64(startTime), float64(endTime))
	if err != nil {
		log.Println("fail to get flags", err)
		t.Fail()
	}

	log.Println(values)
	log.Println(warnings)
}

// Test Snapshot.
func TestSnapshot(t *testing.T) {
	fmt.Println("Get a snapshot")
	snapshot, err := client.Snapshot("test", true)
	if err != nil {
		log.Println("fail to get snapshot", err)
		t.Fail()
	}

	log.Println(snapshot)
}

func TestRules(t *testing.T) {
	fmt.Println("Get Rules")
	rules, err := client.Rules("test")
	if err != nil {
		log.Println("fail to get rules", err)
		t.Fail()
	}

	log.Println(rules)
}

func TestTargets(t *testing.T) {
	fmt.Println("Get Targets")
	targets, err := client.Targets("test")
	if err != nil {
		log.Println("fail to get targets", err)
		t.Fail()
	}

	log.Println(targets)
}

func TestTargetsMetadata(t *testing.T) {
	fmt.Println("Get Targets Metadata")
	rules, err := client.TargetsMetadata("test", "{job=\"prometheus\"}", "go_goroutines", "1")
	if err != nil {
		log.Println("fail to get targets", err)
		t.Fail()
	}

	log.Println(rules)
}
*/
