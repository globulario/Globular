package cluster

import (
	"fmt"
	"strings"
	"testing"
)

// joinScriptForTest returns the raw joinScript template expanded with fixed
// test arguments. It does NOT call the HTTP handler (which reads real CA files)
// so that tests run without cluster infrastructure.
func joinScriptForTest() string {
	return fmt.Sprintf(joinScript,
		"FAKECAB64==",    // %[1]s CA cert
		"10.0.0.1:8443",  // %[2]s gateway
		"10.0.0.1:12000", // %[3]s controller
		"",               // %[4]s scylla GPG (empty)
		"",               // %[5]s scylla apt source (empty)
	)
}

// TestJoinScript_NoLocalMinioHostsEntry verifies that the join script does NOT
// write the joining node's own IP as a minio.globular.internal resolver.
//
// A Day-1 joining node must not assume it is a MinIO/objectstore member.
// Objectstore membership is controlled only by ObjectStoreDesiredState in etcd
// and the topology workflow — never by the join script.
func TestJoinScript_NoLocalMinioHostsEntry(t *testing.T) {
	script := joinScriptForTest()

	forbidden := `"${NODE_IP} minio.globular.internal`
	if strings.Contains(script, forbidden) {
		t.Errorf("join script must not write NODE_IP as minio.globular.internal resolver\n"+
			"found: %q\n"+
			"reason: objectstore membership is governed by topology contract, not join script",
			forbidden)
	}
}

// TestJoinScript_BootstrapMinioFallbackPresent verifies that the join script
// writes the bootstrap node as the minio.globular.internal fallback.
//
// Before DNS and the full cluster stack are running, artifact fetches must
// still be able to reach MinIO. The bootstrap node is the existing objectstore
// endpoint and is the correct pre-DNS fallback.
func TestJoinScript_BootstrapMinioFallbackPresent(t *testing.T) {
	script := joinScriptForTest()

	// The bootstrap fallback line must be present.
	want := `"${BOOTSTRAP_HOST} minio.globular.internal`
	if !strings.Contains(script, want) {
		t.Errorf("join script must write BOOTSTRAP_HOST as minio.globular.internal fallback\n"+
			"expected to find: %q", want)
	}
}

// TestJoinScript_TopologyContractComment verifies that the join script contains
// a comment explaining the topology contract and why this node is not admitted
// as a MinIO member automatically.
func TestJoinScript_TopologyContractComment(t *testing.T) {
	script := joinScriptForTest()

	markers := []string{
		"topology contract",
		"ObjectStoreDesiredState",
		"apply-topology",
	}
	for _, m := range markers {
		if !strings.Contains(script, m) {
			t.Errorf("join script missing topology contract documentation: expected to find %q", m)
		}
	}
}

// TestJoinScript_MinioHostsExactlyOneLine verifies that exactly one
// minio.globular.internal /etc/hosts entry is written (the bootstrap fallback).
func TestJoinScript_MinioHostsExactlyOneLine(t *testing.T) {
	script := joinScriptForTest()

	var hostsLines []string
	for _, line := range strings.Split(script, "\n") {
		if strings.Contains(line, "minio.globular.internal") &&
			strings.HasPrefix(strings.TrimSpace(line), "echo") &&
			strings.Contains(line, "/etc/hosts") {
			hostsLines = append(hostsLines, strings.TrimSpace(line))
		}
	}

	if len(hostsLines) != 1 {
		t.Errorf("expected exactly 1 minio.globular.internal /etc/hosts echo line, got %d:\n%s",
			len(hostsLines), strings.Join(hostsLines, "\n"))
	}
}
