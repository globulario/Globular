package cluster

import (
	"fmt"
	"os"
	"os/exec"
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

// ─── MinIO / objectstore topology contract ───────────────────────────────────

// TestJoinScript_NoLocalMinioHostsEntry verifies that the join script does NOT
// write the joining node's own IP as a minio.globular.internal resolver.
//
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
func TestJoinScript_BootstrapMinioFallbackPresent(t *testing.T) {
	script := joinScriptForTest()

	want := `"${BOOTSTRAP_HOST} minio.globular.internal`
	if !strings.Contains(script, want) {
		t.Errorf("join script must write BOOTSTRAP_HOST as minio.globular.internal fallback\n"+
			"expected to find: %q", want)
	}
}

// TestJoinScript_TopologyContractComment verifies that the join script contains
// the topology contract explanation.
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

// TestJoinScript_NoMinioServiceStart verifies that the join script never starts
// globular-minio.service. Objectstore topology is exclusively controller-driven.
func TestJoinScript_NoMinioServiceStart(t *testing.T) {
	script := joinScriptForTest()

	forbidden := []string{
		"systemctl start globular-minio.service",
		"systemctl enable globular-minio.service\nsystemctl start",
		"systemctl enable --now globular-minio.service",
	}
	for _, f := range forbidden {
		if strings.Contains(script, f) {
			t.Errorf("join script must not start globular-minio.service\n"+
				"found: %q\n"+
				"reason: objectstore topology is exclusively governed by ObjectStoreDesiredState",
				f)
		}
	}
}

// ─── etcd join invariants ────────────────────────────────────────────────────

// TestJoinScript_EtcdYamlExistingClusterState verifies that the join script
// always writes initial-cluster-state: "existing". Writing "new" on a Day-1
// joining node would fork the etcd cluster and destroy quorum.
func TestJoinScript_EtcdYamlExistingClusterState(t *testing.T) {
	script := joinScriptForTest()

	want := `initial-cluster-state: "existing"`
	if !strings.Contains(script, want) {
		t.Errorf("join script must write initial-cluster-state: \"existing\" in etcd.yaml\n"+
			"expected to find: %q", want)
	}
}

// TestJoinScript_NoSingleNodeEtcdSeed verifies that the join script never
// writes initial-cluster-state: "new". Only the Day-0 bootstrap uses "new".
func TestJoinScript_NoSingleNodeEtcdSeed(t *testing.T) {
	script := joinScriptForTest()

	forbidden := `initial-cluster-state: "new"`
	if strings.Contains(script, forbidden) {
		t.Errorf("join script must not write initial-cluster-state: \"new\"\n"+
			"found: %q\n"+
			"reason: \"new\" is only for Day-0 single-node bootstrap; Day-1 join must use \"existing\"",
			forbidden)
	}
}

// TestJoinScript_GhostMemberRemovalBeforeMemberAdd verifies that ghost/stale
// etcd member removal appears before the member add call in the script.
// This is required to unblock member add when a prior failed join left a
// registered peer URL without completing etcd startup.
func TestJoinScript_GhostMemberRemovalBeforeMemberAdd(t *testing.T) {
	script := joinScriptForTest()

	// Match the actual etcdctl invocations, not header comments.
	removeIdx := strings.Index(script, `member remove "${GHOST_ID}"`)
	addIdx := strings.Index(script, `member add "${ETCD_NAME}"`)
	if removeIdx < 0 {
		t.Fatal("join script must contain 'member remove \"${GHOST_ID}\"' for ghost cleanup")
	}
	if addIdx < 0 {
		t.Fatal("join script must contain 'member add \"${ETCD_NAME}\"'")
	}
	if removeIdx > addIdx {
		t.Errorf("'member remove' (ghost cleanup) must appear before 'member add'\n"+
			"remove at byte %d, add at byte %d", removeIdx, addIdx)
	}
}

// TestJoinScript_BackupBeforeWipeInRepairMode verifies that the repair-etcd
// branch backs up the etcd data directory BEFORE wiping it. Wipe without
// backup is not permitted.
func TestJoinScript_BackupBeforeWipeInRepairMode(t *testing.T) {
	script := joinScriptForTest()

	// The backup copy must appear before the rm -rf.
	cpIdx := strings.Index(script, `cp -a "${STATE_DIR}/etcd" "${ETCD_BACKUP}"`)
	rmIdx := strings.Index(script, `rm -rf "${STATE_DIR}/etcd"`)

	if cpIdx < 0 {
		t.Fatal("join script must contain 'cp -a ... ETCD_BACKUP' backup step in repair mode")
	}
	if rmIdx < 0 {
		t.Fatal("join script must contain 'rm -rf ... etcd' wipe step in repair mode")
	}
	if cpIdx > rmIdx {
		t.Errorf("backup (cp -a) must appear before wipe (rm -rf)\n"+
			"cp at byte %d, rm at byte %d", cpIdx, rmIdx)
	}
}

// TestJoinScript_ExistingEtcdDataWithoutRepairFlagFails verifies that the
// script contains a guard that fails if etcd data exists and --repair-etcd
// was not given. This prevents silent WAL destruction on a live node.
func TestJoinScript_ExistingEtcdDataWithoutRepairFlagFails(t *testing.T) {
	script := joinScriptForTest()

	// Must contain the guard that checks for existing etcd data.
	if !strings.Contains(script, "--repair-etcd was not given") {
		t.Error("join script must fail if etcd data exists and --repair-etcd was not given")
	}
}

// TestJoinScript_EtcdFailIsFatal verifies that the script does not silently
// continue after an etcd health check failure. Soft failures like
// "Continuing (controller may recover)" are not acceptable — the node-agent
// must never start before etcd is healthy.
func TestJoinScript_EtcdFailIsFatal(t *testing.T) {
	script := joinScriptForTest()

	softFails := []string{
		"Continuing (controller may recover)",
		"Continuing...",
	}
	for _, s := range softFails {
		if strings.Contains(script, s) {
			t.Errorf("join script must not soft-fail after etcd health check failure\n"+
				"found: %q\n"+
				"reason: node-agent must never start before etcd is healthy", s)
		}
	}

	// The die call on etcd failure must be present.
	if !strings.Contains(script, "etcd join failed") {
		t.Error("join script must contain a fatal error message when etcd health check fails")
	}
}

// TestJoinScript_NoLocalhostPeerURLInEtcdYaml verifies that the etcd.yaml
// written by the join script never uses localhost or 127.0.0.1 in peer URLs.
// Loopback peer URLs would make the member unreachable from other nodes.
func TestJoinScript_NoLocalhostPeerURLInEtcdYaml(t *testing.T) {
	script := joinScriptForTest()

	// Find the etcd.yaml heredoc section (between <<ETCDCFG and ETCDCFG).
	startMark := `> "${STATE_DIR}/config/etcd.yaml" <<ETCDCFG`
	endMark := "ETCDCFG"
	start := strings.Index(script, startMark)
	if start < 0 {
		t.Fatal("could not find etcd.yaml heredoc start marker")
	}
	after := script[start+len(startMark):]
	end := strings.Index(after, "\n"+endMark)
	if end < 0 {
		t.Fatal("could not find etcd.yaml heredoc end marker")
	}
	etcdYaml := after[:end]

	forbidden := []string{
		"127.0.0.1:2380",
		"localhost:2380",
		"[::1]:2380",
	}
	for _, f := range forbidden {
		if strings.Contains(etcdYaml, f) {
			t.Errorf("etcd.yaml must not contain loopback peer URL %q\n"+
				"reason: loopback peer URLs make this member unreachable from other nodes",
				f)
		}
	}
}

// TestJoinScript_LoopbackPeerNormalization verifies that the join script
// normalizes loopback addresses in existing member peer URLs to the bootstrap
// host's real IP. The Day-0 founding node may have registered with 127.0.0.1.
func TestJoinScript_LoopbackPeerNormalization(t *testing.T) {
	script := joinScriptForTest()

	// The sed command that normalizes loopback peer URLs must be present.
	if !strings.Contains(script, `127\\.0\\.0\\.1`) || !strings.Contains(script, "BOOTSTRAP_HOST") {
		t.Error("join script must normalize loopback peer URLs in initial-cluster to BOOTSTRAP_HOST")
	}
}

// TestJoinScript_NodeAgentStartAfterEtcdHealthGate verifies that the
// node-agent start command appears AFTER the etcd health verification loop
// in the script (not just via systemd ordering).
func TestJoinScript_NodeAgentStartAfterEtcdHealthGate(t *testing.T) {
	script := joinScriptForTest()

	healthGateIdx := strings.Index(script, `etcd did not become healthy`)
	agentStartIdx := strings.Index(script, `systemctl start globular-node-agent.service`)

	if healthGateIdx < 0 {
		t.Fatal("join script must contain etcd health gate failure message")
	}
	if agentStartIdx < 0 {
		t.Fatal("join script must contain node-agent start command")
	}
	if agentStartIdx < healthGateIdx {
		t.Errorf("node-agent start (byte %d) must appear after etcd health gate (byte %d)\n"+
			"reason: node-agent must never start before etcd is verified healthy",
			agentStartIdx, healthGateIdx)
	}
}

// ─── Service start ordering ──────────────────────────────────────────────────

// TestJoinScript_NodeAgentAfterEtcd verifies that the node-agent systemd unit
// declares After=...globular-etcd.service so systemd enforces the dependency.
func TestJoinScript_NodeAgentAfterEtcd(t *testing.T) {
	script := joinScriptForTest()

	want := "After=network-online.target globular-etcd.service"
	if !strings.Contains(script, want) {
		t.Errorf("node-agent unit must declare After=...globular-etcd.service\n"+
			"expected to find: %q", want)
	}
}

// TestJoinScript_NodeAgentRequiresEtcd verifies that the node-agent systemd
// unit declares Requires=globular-etcd.service.
func TestJoinScript_NodeAgentRequiresEtcd(t *testing.T) {
	script := joinScriptForTest()

	want := "Requires=globular-etcd.service"
	if !strings.Contains(script, want) {
		t.Errorf("node-agent unit must declare Requires=globular-etcd.service\n"+
			"expected to find: %q", want)
	}
}

// ─── Malformed command fragments ─────────────────────────────────────────────

// TestJoinScript_NoStatusExitFragment verifies that the join script does not
// contain "systemctl status <service> exit" — a class of bug where 'exit' is
// accidentally concatenated with a status check as an argument instead of
// being a separate command on the next line.
func TestJoinScript_NoStatusExitFragment(t *testing.T) {
	script := joinScriptForTest()

	lines := strings.Split(script, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Match: systemctl status <anything> exit
		// This pattern passes "exit" as a flag to systemctl, not as a shell command.
		if strings.HasPrefix(trimmed, "systemctl status") && strings.HasSuffix(trimmed, " exit") {
			t.Errorf("malformed command on line %d: %q\n"+
				"'exit' must be on its own line, not appended to systemctl status",
				i+1, trimmed)
		}
		// Also catch: systemctl status <service> exit <code>
		if strings.HasPrefix(trimmed, "systemctl status") {
			fields := strings.Fields(trimmed)
			for j, f := range fields {
				if f == "exit" && j > 1 {
					t.Errorf("malformed command on line %d: %q\n"+
						"'exit' appears as an argument to systemctl status",
						i+1, trimmed)
				}
			}
		}
	}
}

// TestJoinScript_PipefailSet verifies the script uses 'set -euo pipefail'
// (not just 'set -eu'), ensuring pipeline failures are caught.
func TestJoinScript_PipefailSet(t *testing.T) {
	script := joinScriptForTest()

	if !strings.Contains(script, "set -euo pipefail") {
		t.Error("join script must use 'set -euo pipefail' (not just 'set -eu')")
	}
}

// TestJoinScript_TargetedServiceStop verifies the join script stops only
// named services (globular-etcd, globular-node-agent) rather than using
// a wildcard pattern like 'systemctl stop globular-*.service'.
func TestJoinScript_TargetedServiceStop(t *testing.T) {
	script := joinScriptForTest()

	forbidden := []string{
		"systemctl stop 'globular-*.service'",
		`systemctl stop "globular-*.service"`,
		"systemctl stop globular-*.service",
		"systemctl disable 'globular-*.service'",
	}
	for _, f := range forbidden {
		if strings.Contains(script, f) {
			t.Errorf("join script must not wildcard-stop all globular services\n"+
				"found: %q\n"+
				"reason: other globular services (xDS, envoy, etc.) must not be disrupted",
				f)
		}
	}
}

// TestJoinScript_Phase57UsesBootstrapEtcd verifies that the Phase 5.7 member
// presence check queries BOOTSTRAP_ETCD, not the local NODE_IP:2379 endpoint.
// A locally-healthy etcd that forked its own cluster passes a local check;
// only the bootstrap view proves the node actually joined the existing cluster.
func TestJoinScript_Phase57UsesBootstrapEtcd(t *testing.T) {
	script := joinScriptForTest()

	// Find Phase 5.7 block.
	marker := "[5.7] Verifying named member presence"
	start := strings.Index(script, marker)
	if start < 0 {
		t.Fatal("join script must contain Phase 5.7 member verification block")
	}
	// Extract enough context to cover the etcdctl invocation.
	block := script[start : start+400]

	if strings.Contains(block, `"https://${NODE_IP}:2379"`) {
		t.Error("Phase 5.7 must NOT query local NODE_IP:2379 — a forked standalone etcd passes that check\n" +
			"must query BOOTSTRAP_ETCD instead")
	}
	if !strings.Contains(block, `"${BOOTSTRAP_ETCD}"`) {
		t.Error("Phase 5.7 must query BOOTSTRAP_ETCD to verify the node is visible in the existing cluster")
	}
}

// TestJoinScript_Phase57DiesOnMissingMember verifies that Phase 5.7 calls die
// (not log_warn) when the joining node is not visible in the bootstrap cluster.
func TestJoinScript_Phase57DiesOnMissingMember(t *testing.T) {
	script := joinScriptForTest()

	if strings.Contains(script, "not yet visible — may need propagation time") {
		t.Error("Phase 5.7 must not silently warn when member is not visible — must die")
	}
	if !strings.Contains(script, "forked its own cluster") {
		t.Error("Phase 5.7 die message must mention forked cluster to guide operator remediation")
	}
}

// ─── bash syntax validation ──────────────────────────────────────────────────

// TestJoinScript_BashNSyntaxCheck writes the expanded script to a temp file
// and runs 'bash -n' to detect syntax errors. Skipped if bash is not available.
func TestJoinScript_BashNSyntaxCheck(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not found in PATH — skipping syntax check")
	}

	script := joinScriptForTest()

	f, err := os.CreateTemp("", "globular-join-*.sh")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(script); err != nil {
		t.Fatalf("writing script: %v", err)
	}
	f.Close()

	out, err := exec.Command("bash", "-n", f.Name()).CombinedOutput()
	if err != nil {
		t.Errorf("join script has bash syntax errors:\n%s", string(out))
	}
}

// TestJoinScript_ShellcheckIfAvailable runs shellcheck on the expanded script
// if shellcheck is installed. Not required for CI but catches common mistakes.
func TestJoinScript_ShellcheckIfAvailable(t *testing.T) {
	if _, err := exec.LookPath("shellcheck"); err != nil {
		t.Skip("shellcheck not found in PATH — skipping (install for stronger linting)")
	}

	script := joinScriptForTest()

	f, err := os.CreateTemp("", "globular-join-*.sh")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(script); err != nil {
		t.Fatalf("writing script: %v", err)
	}
	f.Close()

	// SC2034: unused variables (template vars like SCYLLA_APT_SOURCE are used
	//         conditionally so shellcheck flags them as unused).
	// SC2086: double-quote prevention (some intentional word-splitting in awk).
	out, err := exec.Command(
		"shellcheck",
		"--exclude=SC2034,SC2086",
		f.Name(),
	).CombinedOutput()
	if err != nil {
		t.Errorf("shellcheck found issues in join script:\n%s", string(out))
	}
}
