package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	servicesconfig "github.com/globulario/services/golang/config"
)

// joinScriptForTest returns the raw joinScript template expanded with fixed
// test arguments. It does NOT call the HTTP handler (which reads real CA files)
// so that tests run without cluster infrastructure.
func joinScriptForTest() string {
	return fmt.Sprintf(joinScript,
		"FAKECAB64==",                   // %[1]s CA cert
		"10.0.0.1:8443",                 // %[2]s gateway
		"10.0.0.1:12000",                // %[3]s controller
		"1.2.0",                         // %[4]s platform version
		servicesconfig.EtcdClusterToken, // %[5]s etcd initial-cluster-token
	)
}

// TestJoinScript_EtcdInitialClusterTokenPresent verifies the Day-1 join writes a
// complete etcd.yaml: the persisted config must carry initial-cluster-token set
// to the fixed config.EtcdClusterToken constant. Without it the joined node's
// etcd.yaml is incomplete — tolerated only because a normally-joined member
// reads membership from its data-dir and ignores initial-* flags, but unsafe on
// a data-dir-loss re-bootstrap/repair. (Founding node renders this; Day-1 join
// previously omitted it.)
func TestJoinScript_EtcdInitialClusterTokenPresent(t *testing.T) {
	script := joinScriptForTest()

	if !strings.Contains(script, `initial-cluster-token: "${ETCD_CLUSTER_TOKEN}"`) {
		t.Error("etcd.yaml heredoc must write initial-cluster-token from ${ETCD_CLUSTER_TOKEN}")
	}
	want := fmt.Sprintf(`ETCD_CLUSTER_TOKEN="%s"`, servicesconfig.EtcdClusterToken)
	if !strings.Contains(script, want) {
		t.Errorf("join script must set ETCD_CLUSTER_TOKEN to config.EtcdClusterToken (%q); not found in rendered script",
			servicesconfig.EtcdClusterToken)
	}
}

// ─── platform version resolution (Day-1 bootstrap source) ────────────────────

// TestResolvePlatformVersion_AnchorWins verifies the etcd active_release anchor
// takes precedence over the local release-index.json projection and the binary
// version. Day-1 must bootstrap the version the cluster is actually running, not
// a stale on-disk projection — the root cause of a joining node installing an
// old node-agent that never converged.
func TestResolvePlatformVersion_AnchorWins(t *testing.T) {
	anchor := func() string { return "1.2.251" }
	readFile := func(string) ([]byte, error) {
		return []byte(`{"platform_release":"1.2.240"}`), nil // stale projection
	}
	if got := resolvePlatformVersion("0.0.0-dev", anchor, readFile); got != "1.2.251" {
		t.Errorf("anchor must win over projection and binary: got %q, want 1.2.251", got)
	}
}

// TestResolvePlatformVersion_FallsBackToProjection verifies that when the anchor
// is empty (etcd unreachable / key absent), the on-disk projection is used.
func TestResolvePlatformVersion_FallsBackToProjection(t *testing.T) {
	anchor := func() string { return "" }
	readFile := func(string) ([]byte, error) {
		return []byte(`{"platform_release":"1.2.240"}`), nil
	}
	if got := resolvePlatformVersion("0.0.0-dev", anchor, readFile); got != "1.2.240" {
		t.Errorf("must fall back to projection when anchor empty: got %q, want 1.2.240", got)
	}
}

// TestResolvePlatformVersion_FallsBackToBinary verifies that with no anchor and
// no readable projection, the gateway binary's own version is used (Day-0).
func TestResolvePlatformVersion_FallsBackToBinary(t *testing.T) {
	anchor := func() string { return "" }
	readFile := func(string) ([]byte, error) { return nil, os.ErrNotExist }
	if got := resolvePlatformVersion("1.2.99", anchor, readFile); got != "1.2.99" {
		t.Errorf("must fall back to binary version (Day-0): got %q, want 1.2.99", got)
	}
}

// TestResolvePlatformVersion_EmptyProjectionFallsThrough verifies that a
// projection with an empty platform_release does not shadow the binary version.
func TestResolvePlatformVersion_EmptyProjectionFallsThrough(t *testing.T) {
	anchor := func() string { return "" }
	readFile := func(string) ([]byte, error) { return []byte(`{"platform_release":""}`), nil }
	if got := resolvePlatformVersion("1.2.99", anchor, readFile); got != "1.2.99" {
		t.Errorf("empty projection must fall through to binary: got %q, want 1.2.99", got)
	}
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

// TestJoinScript_EtcdJoinsAsLearner locks in the learner-first join contract
// (task#17 c): a Day-1 node MUST join etcd as a non-voting learner so a slow or
// failed join can never break the existing cluster's quorum. Regressing to a
// full-voter member add reintroduces the 1->2 quorum trap.
func TestJoinScript_EtcdJoinsAsLearner(t *testing.T) {
	script := joinScriptForTest()

	// The member add must carry --learner.
	if !strings.Contains(script, `member add "${ETCD_NAME}" --learner`) {
		t.Fatal("join script must add the node as a non-voting learner (member add ... --learner)")
	}
	// The max-learners=1 pre-add wait must be present (etcd allows one learner at a time).
	if !strings.Contains(script, `"isLearner":true`) {
		t.Error("join script must wait for a prior learner to be promoted before adding (max-learners=1)")
	}
	// The local readiness gate command must be `endpoint status` (learner-safe),
	// NOT `endpoint health` (a learner rejects the read RPC it needs). Match the
	// actual command invocation, not explanatory comment text.
	if !strings.Contains(script, "endpoint status >/dev/null") {
		t.Error("learner join must gate local readiness on 'endpoint status' (learners reject 'endpoint health')")
	}
	if strings.Contains(script, "endpoint health >/dev/null") {
		t.Error("join script must not run 'endpoint health' against the local learner — it is rejected with 'rpc not supported for learner'")
	}
}

// TestJoinScript_LearnerMemberAddRetriesWithinBudget locks the contract broken on
// 2026-08-01: the learner-slot budget must cover the member ADD, not only the
// wait that precedes it.
//
// etcd allows one learner (max-learners=1) and serializes raft configuration
// changes, so observing a free slot does not reserve it. Joining four nodes
// concurrently, node-4 and node-5 both left the wait, fired a single `member add`
// into the contended slot, and died at 30s on "etcdserver: request timed out" —
// having spent none of their 300s budget, because the add sat outside the loop.
//
// A single-attempt add is therefore forbidden here regardless of how generous the
// preceding wait is.
func TestJoinScript_LearnerMemberAddRetriesWithinBudget(t *testing.T) {
	script := joinScriptForTest()

	addIdx := strings.Index(script, `member add "${ETCD_NAME}" --learner`)
	if addIdx < 0 {
		t.Fatal("join script must contain the learner member add")
	}

	// The add must sit inside the budgeted loop, i.e. the loop must still be open
	// where the add appears. Anchor on the deadline the loop is bounded by.
	deadlineIdx := strings.Index(script, "ETCD_ADD_DEADLINE=")
	if deadlineIdx < 0 {
		t.Fatal("learner member add must be bounded by a wall-clock budget (ETCD_ADD_DEADLINE)")
	}
	if deadlineIdx > addIdx {
		t.Error("the learner-slot budget must be established BEFORE the member add,\n" +
			"otherwise the add is a single attempt outside the budget (2026-08-01 node-4/node-5 failure)")
	}

	// Contention must be classified as retryable. "too many learner members" is
	// the max-learners rejection; "request timed out" is the blocked raft config
	// change that actually killed node-4 and node-5.
	for _, retryable := range []string{"too many learner members", "request timed out"} {
		if !strings.Contains(script, retryable) {
			t.Errorf("join script must treat %q as retryable contention on the learner slot,\n"+
				"not as a terminal rejection", retryable)
		}
	}

	// Each attempt must be bounded so one blocked add cannot consume the whole
	// budget in a single call.
	if !strings.Contains(script, "--command-timeout=") {
		t.Error("each learner member add attempt must be time-bounded (--command-timeout),\n" +
			"or a blocked add swallows the entire retry budget")
	}

	// Exhausting the budget must remain fatal — retrying must not become a way to
	// proceed without membership.
	if !strings.Contains(script, "did not succeed within 300s") {
		t.Error("budget exhaustion on member add must be fatal, not a soft continue")
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

// TestJoinScript_TarballDownloadFallsBackToGateway verifies that the generated
// join script tries GitHub first, then falls back to the controller-local
// bundle endpoint when the platform tag is not published upstream.
func TestJoinScript_TarballDownloadFallsBackToGateway(t *testing.T) {
	script := joinScriptForTest()

	if !strings.Contains(script, `GITHUB_TARBALL_URL="https://github.com/globulario/services/releases/download/v${PLATFORM_VERSION}/${TARBALL_NAME}"`) {
		t.Fatal("join script must define the GitHub tarball URL")
	}
	if !strings.Contains(script, `GATEWAY_TARBALL_URL="https://${GATEWAY}/join/bundle/${TARBALL_NAME}"`) {
		t.Fatal("join script must define the gateway fallback tarball URL")
	}
	if !strings.Contains(script, `GitHub release tarball unavailable for v${PLATFORM_VERSION} — falling back to controller bundle`) {
		t.Fatal("join script must log the controller fallback when GitHub is unavailable")
	}
	if !strings.Contains(script, `could not download ${TARBALL_NAME} (+.sha256) from GitHub or controller fallback`) {
		t.Fatal("join script must fail closed if both GitHub and controller bundle download fail")
	}
	// The controller fallback must fetch and verify a checksum, not trust the
	// bytes blindly (mirrors the GitHub path).
	if !strings.Contains(script, `${GATEWAY_TARBALL_URL}.sha256`) {
		t.Fatal("join script must fetch the controller fallback bundle .sha256")
	}
	if !strings.Contains(script, `controller fallback bundle checksum mismatch`) {
		t.Fatal("join script must fail closed on a controller fallback checksum mismatch")
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

	// The parameter expansion that normalizes loopback peer URLs must be present.
	if !strings.Contains(script, "127.0.0.1/${BOOTSTRAP_HOST}") {
		t.Error("join script must normalize loopback peer URLs in initial-cluster to BOOTSTRAP_HOST")
	}
}

// TestJoinScript_NodeAgentStartAfterEtcdHealthGate verifies that the
// node-agent start command appears AFTER the etcd health verification loop
// in the script (not just via systemd ordering).
func TestJoinScript_NodeAgentStartAfterEtcdHealthGate(t *testing.T) {
	script := joinScriptForTest()

	// Learner-first join: a learner cannot serve client health RPCs, so the local
	// readiness gate uses `endpoint status` and dies with "etcd did not come up".
	healthGateIdx := strings.Index(script, `etcd did not come up`)
	agentStartIdx := strings.Index(script, `systemctl start globular-node-agent.service`)

	if healthGateIdx < 0 {
		t.Fatal("join script must contain etcd readiness gate failure message")
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

// TestJoinScript_Day1InfraControlPlaneOrder verifies the requested Day-1
// bootstrap ladder: etcd cluster membership is proven first, then local
// ScyllaDB, then the local control-plane stack, then node-agent admission.
func TestJoinScript_Day1InfraControlPlaneOrder(t *testing.T) {
	script := joinScriptForTest()

	markers := []string{
		`member ${ETCD_NAME} confirmed visible in cluster`,
		`log_phase "6 — ScyllaDB Local Join"`,
		`log_phase "7 — Local Control Plane"`,
		`log_phase "8 — node-agent"`,
		`log_phase "9 — MinIO / Objectstore Hold Enforcement"`,
	}
	last := -1
	for _, marker := range markers {
		idx := strings.Index(script, marker)
		if idx < 0 {
			t.Fatalf("join script missing marker %q", marker)
		}
		if idx <= last {
			t.Fatalf("join script marker %q is out of order", marker)
		}
		last = idx
	}
}

// TestJoinScript_LocalBootstrapUsesInstallerPackages verifies the Day-1 local
// bootstrap path goes through globular-installer and staged distro packages.
func TestJoinScript_LocalBootstrapUsesInstallerPackages(t *testing.T) {
	script := joinScriptForTest()

	for _, marker := range []string{
		`find "${PACKAGES_DIR}" -maxdepth 1 -name "${name}_*_linux_amd64.tgz"`,
		`"${INSTALL_DIR}/globular-installer" install`,
		`install_pkg "scylladb"`,
		`install_pkg "xds"`,
		`install_pkg "envoy"`,
		`install_pkg "gateway"`,
		`install_pkg "workflow"`,
		`install_pkg "repository"`,
		`install_pkg "cluster-controller"`,
	} {
		if !strings.Contains(script, marker) {
			t.Fatalf("join script must contain installer/package marker %q", marker)
		}
	}
}

// TestJoinScript_LocalControlPlaneOrder verifies dependency-sensitive local
// service order after ScyllaDB: xDS before Envoy, Envoy before Gateway, then
// workflow/repository before the local cluster-controller.
func TestJoinScript_LocalControlPlaneOrder(t *testing.T) {
	script := joinScriptForTest()

	markers := []string{
		`install_pkg "scylladb"`,
		`install_pkg "xds"`,
		`install_pkg "envoy"`,
		`install_pkg "gateway"`,
		`install_pkg "workflow"`,
		`install_pkg "repository"`,
		`install_pkg "cluster-controller"`,
		`systemctl start globular-node-agent.service`,
	}
	last := -1
	for _, marker := range markers {
		idx := strings.Index(script, marker)
		if idx < 0 {
			t.Fatalf("join script missing marker %q", marker)
		}
		if idx <= last {
			t.Fatalf("join script marker %q is out of order", marker)
		}
		last = idx
	}
}

// TestJoinScript_ScyllaReadinessBeforeControlPlane verifies that local
// control-plane services do not start before ScyllaDB's unit and CQL port are
// both ready.
func TestJoinScript_ScyllaReadinessBeforeControlPlane(t *testing.T) {
	script := joinScriptForTest()

	unitIdx := strings.Index(script, `wait_unit_active "scylla-server.service" 300 "ScyllaDB"`)
	portIdx := strings.Index(script, `wait_tcp_port "${NODE_IP}" 9042 300 "ScyllaDB CQL"`)
	xdsIdx := strings.Index(script, `install_pkg "xds"`)
	if unitIdx < 0 || portIdx < 0 || xdsIdx < 0 {
		t.Fatalf("join script missing Scylla readiness or xDS install marker")
	}
	if !(unitIdx < portIdx && portIdx < xdsIdx) {
		t.Fatalf("Scylla unit and CQL readiness must be verified before xDS install")
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

// ─── Preflight lifecycle contract (scar 2026-07-10) ──────────────────────────
//
// Re-running the fresh-join script on an already-joined node used to stop its
// etcd member and node-agent DURING preflight, then refuse on the existing
// data directory — leaving the member down and collapsing quorum on a 2-node
// cluster. The contract is now: every refusal fires BEFORE any service stop,
// and an armed compensation trap restores previously-active units on abort.

// TestJoinScript_LifecycleRefusalBeforeServiceStop verifies that both refusal
// paths (non-pristine node without --repair-etcd; repair mode against a live
// member) appear before the first service stop in the script.
func TestJoinScript_LifecycleRefusalBeforeServiceStop(t *testing.T) {
	script := joinScriptForTest()

	stopIdx := strings.Index(script, `systemctl stop "${unit}"`)
	if stopIdx < 0 {
		t.Fatal("join script must contain the targeted service-stop loop")
	}

	classifyIdx := strings.Index(script, "[0.0] Classifying node lifecycle state")
	if classifyIdx < 0 {
		t.Fatal("join script must classify node lifecycle state in step [0.0]")
	}
	if classifyIdx > stopIdx {
		t.Errorf("lifecycle classification (byte %d) must run before the service-stop loop (byte %d)",
			classifyIdx, stopIdx)
	}

	refusals := []string{
		`REFUSED: node is ${LIFECYCLE}`,
		"--repair-etcd was not given",
		"--repair-etcd given but globular-etcd.service is ACTIVE",
	}
	for _, r := range refusals {
		idx := strings.Index(script, r)
		if idx < 0 {
			t.Errorf("join script missing refusal %q", r)
			continue
		}
		if idx > stopIdx {
			t.Errorf("refusal %q (byte %d) must appear BEFORE the service-stop loop (byte %d)\n"+
				"reason: an aborted preflight must leave the node exactly as it found it",
				r, idx, stopIdx)
		}
	}
}

// TestJoinScript_CompensationTrapArmedBeforeStops verifies the abort
// compensation trap is installed and armed before the first service stop, so
// any later abort restores the units that were active when the script started.
func TestJoinScript_CompensationTrapArmedBeforeStops(t *testing.T) {
	script := joinScriptForTest()

	stopIdx := strings.Index(script, `systemctl stop "${unit}"`)
	trapIdx := strings.Index(script, "trap compensate_on_abort EXIT")
	armIdx := strings.Index(script, "COMPENSATION_ARMED=1")

	if stopIdx < 0 || trapIdx < 0 || armIdx < 0 {
		t.Fatalf("join script missing compensation markers (stop=%d trap=%d arm=%d)",
			stopIdx, trapIdx, armIdx)
	}
	if trapIdx > stopIdx || armIdx > stopIdx {
		t.Errorf("compensation trap (byte %d) and arming (byte %d) must precede the service-stop loop (byte %d)",
			trapIdx, armIdx, stopIdx)
	}
	// The trap must restore exactly the units recorded as active at start.
	for _, m := range []string{"ETCD_WAS_ACTIVE", "AGENT_WAS_ACTIVE", "join_interrupted.json"} {
		if !strings.Contains(script, m) {
			t.Errorf("compensation trap missing marker %q", m)
		}
	}
}

// TestJoinScript_UnitRestoreDisarmedAtPointsOfNoReturn verifies unit restore
// is disarmed (a) after the repair-mode wipe — restarting the old unit on an
// empty data dir could fork a new cluster — and (b) after the Phase 5.7
// anti-fork proof, when the local units belong to the new cluster state.
func TestJoinScript_UnitRestoreDisarmedAtPointsOfNoReturn(t *testing.T) {
	script := joinScriptForTest()

	wipeIdx := strings.Index(script, `rm -rf "${STATE_DIR}/etcd"`)
	antiForkIdx := strings.Index(script, "anti-fork verified")
	if wipeIdx < 0 || antiForkIdx < 0 {
		t.Fatalf("join script missing wipe (%d) or anti-fork (%d) marker", wipeIdx, antiForkIdx)
	}

	disarmCount := strings.Count(script, "RESTORE_UNITS_ON_ABORT=0")
	if disarmCount < 2 {
		t.Fatalf("expected at least 2 RESTORE_UNITS_ON_ABORT=0 disarm points (post-wipe, post-anti-fork), got %d", disarmCount)
	}
	firstDisarm := strings.Index(script, "RESTORE_UNITS_ON_ABORT=0")
	if firstDisarm < wipeIdx {
		t.Errorf("first unit-restore disarm (byte %d) must come after the repair-mode wipe (byte %d)",
			firstDisarm, wipeIdx)
	}
	lastDisarm := strings.LastIndex(script, "RESTORE_UNITS_ON_ABORT=0")
	if lastDisarm < antiForkIdx {
		t.Errorf("a unit-restore disarm must follow the anti-fork proof (disarm byte %d, anti-fork byte %d)",
			lastDisarm, antiForkIdx)
	}
}

// TestJoinScript_TrapDisarmedOnSuccess verifies the success path disarms the
// abort trap and clears any stale interrupted-join receipt.
func TestJoinScript_TrapDisarmedOnSuccess(t *testing.T) {
	script := joinScriptForTest()

	disarmIdx := strings.Index(script, "trap - EXIT")
	successIdx := strings.Index(script, "NODE JOINED SUCCESSFULLY")
	clearIdx := strings.Index(script, `rm -f "${STATE_DIR}/nodeagent/join_interrupted.json"`)

	if disarmIdx < 0 {
		t.Fatal("join script must disarm the abort trap (trap - EXIT) on success")
	}
	if clearIdx < 0 {
		t.Fatal("join script must clear the stale interrupted-join receipt on success")
	}
	if successIdx < 0 {
		t.Fatal("join script missing success banner")
	}
	if disarmIdx > successIdx {
		t.Errorf("trap disarm (byte %d) must precede the success banner (byte %d)", disarmIdx, successIdx)
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

// TestJoinScript_ShellcheckIfAvailable runs shellcheck on the expanded script.
// Required for CI — the join script owns critical etcd operations.
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

// ── JoinPlan-driven join: authorization before credential, profile gates ─────

// The JoinPlan authorization must run BEFORE certificate signing: the signed
// cert is an etcd mTLS credential (client-cert-auth trusts the cluster CA),
// so issuance itself is a cluster-affecting step gated on the join token.
func TestJoinScript_AuthorizeBeforeCertSigning(t *testing.T) {
	script := joinScriptForTest()
	authorizeIdx := strings.Index(script, "/join/authorize")
	signIdx := strings.Index(script, "/sign_ca_certificate")
	if authorizeIdx < 0 {
		t.Fatal("join script must contain /join/authorize call")
	}
	if signIdx < 0 {
		t.Fatal("join script must contain /sign_ca_certificate call")
	}
	if authorizeIdx > signIdx {
		t.Errorf("/join/authorize (byte %d) must appear BEFORE /sign_ca_certificate (byte %d)\n"+
			"reason: the CA-signed cert is a credential; it may only be issued to an authorized joiner",
			authorizeIdx, signIdx)
	}
}

// The certificate-signing request must present the controller-issued join_id —
// the gateway refuses to sign CSRs from unauthorized requesters.
func TestJoinScript_CertSigningPresentsJoinID(t *testing.T) {
	script := joinScriptForTest()
	if !strings.Contains(script, "/sign_ca_certificate?csr=${CSR_B64}&join_id=${JOIN_ID}") {
		t.Error("sign_ca_certificate call must include the join_id credential")
	}
}

// The controller-assigned profiles from the JoinPlan drive local installs.
// ScyllaDB joins the storage ring ONLY when the plan assigns 'storage' —
// mirroring the MinIO topology hold.
func TestJoinScript_ScyllaGatedOnStorageProfile(t *testing.T) {
	script := joinScriptForTest()
	gateIdx := strings.Index(script, `if should_install_role "storage"; then`)
	scyllaIdx := strings.Index(script, `install_pkg "scylladb"`)
	if gateIdx < 0 {
		t.Fatal("join script must gate ScyllaDB install on the 'storage' profile")
	}
	if scyllaIdx < 0 {
		t.Fatal("join script must install scylladb via install_pkg")
	}
	if gateIdx > scyllaIdx {
		t.Error("the storage profile gate must appear before the scylladb install")
	}
}

// Control-plane replicas (workflow/repository/cluster-controller) are
// installed ONLY with the 'control-plane' profile. The mesh data plane
// (xds/envoy/gateway) stays unconditional — every node needs it.
func TestJoinScript_ControlPlaneGatedOnProfile(t *testing.T) {
	script := joinScriptForTest()
	gateIdx := strings.Index(script, `if should_install_role "control-plane"; then`)
	ctrlIdx := strings.Index(script, `install_pkg "cluster-controller"`)
	xdsIdx := strings.Index(script, `install_pkg "xds"`)
	if gateIdx < 0 {
		t.Fatal("join script must gate control-plane services on the 'control-plane' profile")
	}
	if ctrlIdx < 0 {
		t.Fatal("join script must install cluster-controller via install_pkg")
	}
	if gateIdx > ctrlIdx {
		t.Error("the control-plane profile gate must appear before the cluster-controller install")
	}
	if xdsIdx < 0 {
		t.Fatal("join script must install xds via install_pkg")
	}
	if xdsIdx > gateIdx {
		t.Error("xds (mesh data plane) must be installed unconditionally, before the control-plane gate")
	}
}

// Legacy fallback: a plan without assigned_profiles must not brick the join —
// the script warns and installs everything (pre-profile controller upgrades).
func TestJoinScript_LegacyPlanFallsBackToFullInstall(t *testing.T) {
	script := joinScriptForTest()
	if !strings.Contains(script, `[[ -z "${ASSIGNED_PROFILES}" ]] && return 0`) {
		t.Error("should_install_role must fall back to installing everything when the plan carries no profiles")
	}
}
