package monitoring_client

import (
	"context"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/monitoring/monitoringpb"

	//	"github.com/davecourtois/Utility"
	"io"

	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Monitoring Client Service
////////////////////////////////////////////////////////////////////////////////

type monitoring_Client struct {
	cc *grpc.ClientConn
	c  monitoringpb.MonitoringServiceClient

	// The name of the service
	name string

	// The ipv4 address
	addresse string

	// The client domain
	domain string

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewMonitoring_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string, token string) *monitoring_Client {
	client := new(monitoring_Client)

	client.addresse = addresse
	client.domain = domain
	client.name = "Monitoring"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
	client.c = monitoringpb.NewMonitoringServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *monitoring_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *monitoring_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *monitoring_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *monitoring_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *monitoring_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *monitoring_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *monitoring_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *monitoring_Client) GetCaFile() string {
	return self.caFile
}

////////////////// Connections management functions //////////////////////////

// Create a new connection.
func (self *monitoring_Client) CreateConnection(id string, host string, storeType float64, port float64) error {
	rqst := &monitoringpb.CreateConnectionRqst{
		Connection: &monitoringpb.Connection{
			Id:    id,
			Host:  host,
			Port:  int32(port),
			Store: monitoringpb.StoreType(int32(storeType)),
		},
	}

	_, err := self.c.CreateConnection(context.Background(), rqst)

	return err
}

// Delete a connection.
func (self *monitoring_Client) DeleteConnection(id string) error {
	rqst := &monitoringpb.DeleteConnectionRqst{
		Id: id,
	}

	_, err := self.c.DeleteConnection(context.Background(), rqst)

	return err
}

// Config returns the current Prometheus configuration.
func (self *monitoring_Client) Config(connectionId string) (string, error) {
	rqst := &monitoringpb.ConfigRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.Config(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResults(), nil
}

// Alerts returns a list of all active alerts.
func (self *monitoring_Client) Alerts(connectionId string) (string, error) {
	rqst := &monitoringpb.AlertsRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.Alerts(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResults(), nil
}

// AlertManagers returns an overview of the current state of the Prometheus alert manager discovery.
func (self *monitoring_Client) AlertManagers(connectionId string) (string, error) {
	rqst := &monitoringpb.AlertManagersRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.AlertManagers(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResults(), nil
}

// CleanTombstones removes the deleted data from disk and cleans up the existing tombstones.
func (self *monitoring_Client) CleanTombstones(connectionId string) error {
	rqst := &monitoringpb.CleanTombstonesRequest{
		ConnectionId: connectionId,
	}

	_, err := self.c.CleanTombstones(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSeries deletes data for a selection of series in a time range.
func (self *monitoring_Client) DeleteSeries(connectionId string, matches []string, startTime float64, endTime float64) error {
	rqst := &monitoringpb.DeleteSeriesRequest{
		ConnectionId: connectionId,
		Matches:      matches,
		StartTime:    startTime,
		EndTime:      endTime,
	}

	_, err := self.c.DeleteSeries(context.Background(), rqst)
	if err != nil {
		return err
	}
	return nil
}

// Flags returns the flag values that Prometheus was launched with.
func (self *monitoring_Client) Flags(connectionId string) (string, error) {
	rqst := &monitoringpb.FlagsRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.Flags(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResults(), nil
}

// LabelNames returns all the unique label names present in the block in sorted order.
func (self *monitoring_Client) LabelNames(connectionId string) ([]string, string, error) {
	rqst := &monitoringpb.LabelNamesRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.LabelNames(context.Background(), rqst)
	if err != nil {
		return nil, "", err
	}

	return rsp.GetLabels(), rsp.GetWarnings(), nil
}

// LabelValues performs a query for the values of the given label.
func (self *monitoring_Client) LabelValues(connectionId string, label string) (string, string, error) {
	rqst := &monitoringpb.LabelValuesRequest{
		ConnectionId: connectionId,
		Label:        label,
	}

	rsp, err := self.c.LabelValues(context.Background(), rqst)
	if err != nil {
		return "", "", err
	}

	return rsp.GetLabelValues(), rsp.GetWarnings(), nil
}

// Query performs a query for the given time.
func (self *monitoring_Client) Query(connectionId string, query string, ts float64) (string, string, error) {
	rqst := &monitoringpb.QueryRequest{
		ConnectionId: connectionId,
		Query:        query,
		Ts:           ts,
	}

	rsp, err := self.c.Query(context.Background(), rqst)
	if err != nil {
		return "", "", err
	}

	return rsp.GetValue(), rsp.GetWarnings(), nil
}

// QueryRange performs a query for the given range.
func (self *monitoring_Client) QueryRange(connectionId string, query string, startTime float64, endTime float64, step float64) (string, string, error) {
	rqst := &monitoringpb.QueryRangeRequest{
		ConnectionId: connectionId,
		Query:        query,
		StartTime:    startTime,
		EndTime:      endTime,
		Step:         step,
	}

	var value string
	var warning string
	stream, err := self.c.QueryRange(context.Background(), rqst)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			return "", "", err
		}

		// Get the result...
		value += msg.GetValue()
		warning = msg.GetWarnings()
	}

	if err != nil {
		return "", "", err
	}

	return value, warning, nil
}

// Series finds series by label matchers.
func (self *monitoring_Client) Series(connectionId string, matches []string, startTime float64, endTime float64) (string, string, error) {
	rqst := &monitoringpb.SeriesRequest{
		ConnectionId: connectionId,
		Matches:      matches,
		StartTime:    startTime,
		EndTime:      endTime,
	}

	rsp, err := self.c.Series(context.Background(), rqst)
	if err != nil {
		return "", "", err
	}

	return rsp.GetLabelSet(), rsp.GetWarnings(), nil
}

// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
// under the TSDB's data directory and returns the directory as response.
func (self *monitoring_Client) Snapshot(connectionId string, skipHead bool) (string, error) {
	rqst := &monitoringpb.SnapshotRequest{
		ConnectionId: connectionId,
		SkipHead:     skipHead,
	}

	rsp, err := self.c.Snapshot(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// Rules returns a list of alerting and recording rules that are currently loaded.
func (self *monitoring_Client) Rules(connectionId string) (string, error) {
	rqst := &monitoringpb.RulesRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.Rules(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// Targets returns an overview of the current state of the Prometheus target discovery.
func (self *monitoring_Client) Targets(connectionId string) (string, error) {
	rqst := &monitoringpb.TargetsRequest{
		ConnectionId: connectionId,
	}

	rsp, err := self.c.Targets(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// TargetsMetadata returns metadata about metrics currently scraped by the target.
func (self *monitoring_Client) TargetsMetadata(connectionId string, matchTarget string, metric string, limit string) (string, error) {
	rqst := &monitoringpb.TargetsMetadataRequest{
		ConnectionId: connectionId,
		MatchTarget:  matchTarget,
		Metric:       metric,
		Limit:        limit,
	}

	rsp, err := self.c.TargetsMetadata(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}
