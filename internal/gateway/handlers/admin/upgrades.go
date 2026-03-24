package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/globulario/services/golang/plan/planpb"
	"github.com/globulario/services/golang/plan/versionutil"
)

// ── Provider interface ──────────────────────────────────────────────────────

// BundleInfo describes a single available package in the repository.
type BundleInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	BuildNumber int64  `json:"build_number,omitempty"`
	Platform    string `json:"platform"`
	PublisherID string `json:"publisher_id,omitempty"`
	SizeBytes   int64  `json:"size_bytes,omitempty"`
	SHA256      string `json:"sha256,omitempty"`
	Kind        string `json:"kind,omitempty"` // "SERVICE", "APPLICATION", "INFRASTRUCTURE" (empty = SERVICE)
}

// UpgradesProvider gives the handler access to the package repository.
type UpgradesProvider interface {
	AvailableBundles(ctx context.Context) ([]BundleInfo, error)
}

// ── Response types ──────────────────────────────────────────────────────────

// UpgradesStatusResponse is the top-level response for GET /admin/upgrades/status.
type UpgradesStatusResponse struct {
	NowUnixMs        int64                `json:"now_unix_ms"`
	Node             string               `json:"node"`
	Platform         string               `json:"platform"`
	Services         []ServiceUpgradeInfo `json:"services"`
	Summary          UpgradesSummary      `json:"summary"`
	RepositoryStatus string               `json:"repository_status"` // ok | unreachable | empty
}

// ServiceUpgradeInfo describes the upgrade status of a single service.
type ServiceUpgradeInfo struct {
	Name              string `json:"name"`
	DisplayName       string `json:"display_name"`
	Category          string `json:"category"`
	Kind              string `json:"kind,omitempty"` // "SERVICE", "APPLICATION", "INFRASTRUCTURE"
	InstalledVersion  string `json:"installed_version"`
	InstalledBuildNum int64  `json:"installed_build_number,omitempty"`
	LatestVersion     string `json:"latest_version"`
	LatestBuildNum    int64  `json:"latest_build_number,omitempty"`
	UpdateAvailable   bool   `json:"update_available"`
	State             string `json:"state"`
	DerivedStatus     string `json:"derived_status"` // healthy | degraded | critical | unknown
	Port              int    `json:"port"`
}

// UpgradesSummary is a rollup of upgrade availability.
type UpgradesSummary struct {
	TotalInstalled      int `json:"total_installed"`
	UpdatesAvailable    int `json:"updates_available"`
	UpToDate            int `json:"up_to_date"`
	Unknown             int `json:"unknown"`
	InfrastructureCount int `json:"infrastructure_count"`
	ServiceCount        int `json:"service_count"`
	ApplicationCount    int `json:"application_count"`
}

// ── Handler ─────────────────────────────────────────────────────────────────

// NewUpgradesHandler returns a GET-only handler for /admin/upgrades/status.
func NewUpgradesHandler(admin AdminProvider, upgrades UpgradesProvider) http.Handler {
	prom := newPromClient("http://localhost:9090", 8*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		platform := runtime.GOOS + "_" + runtime.GOARCH

		// Fetch installed services, available bundles, and Prometheus metrics concurrently.
		var (
			cfgs          []map[string]any
			cfgErr        error
			bundles       []BundleInfo
			bndErr        error
			promMetrics   map[string]*svcMetrics
			promConnected bool
		)

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			cfgs, cfgErr = admin.AllServiceConfigs()
		}()
		go func() {
			defer wg.Done()
			bundles, bndErr = upgrades.AvailableBundles(ctx)
		}()
		go func() {
			defer wg.Done()
			promMetrics, promConnected = fetchPromMetrics(ctx, prom)
		}()
		wg.Wait()

		if cfgErr != nil {
			http.Error(w, "failed to list services: "+cfgErr.Error(), http.StatusInternalServerError)
			return
		}

		// Build latest-version index: baseName -> latest version for this platform.
		latestByName := buildLatestIndex(bundles, platform)

		// Build response.
		var services []ServiceUpgradeInfo
		summary := UpgradesSummary{}

		for _, cfg := range cfgs {
			name := mapStr(cfg, "Name")
			if name == "" {
				continue
			}

			base := normalizeServiceName(name)
			installed := mapStr(cfg, "Version")
			state := mapStr(cfg, "State")
			port := mapInt(cfg, "Port")

			// Cross-validate state with Prometheus (same logic as services handler).
			var pm *svcMetrics
			if promMetrics != nil {
				pm = promMetrics[base]
			}
			derivedStatus, _ := deriveServiceHealth(state, pm, promConnected, port)

			// Correct stale etcd state: if Prometheus proves the process is
			// alive but etcd says stopped/closed, show the true state.
			lowerState := strings.ToLower(state)
			if promConnected && pm != nil && (lowerState == "stopped" || lowerState == "closed") {
				state = "running"
			}

			// Map derived status to a user-friendly state if etcd state is unhelpful.
			if state == "" || strings.ToLower(state) == "unknown" {
				switch derivedStatus {
				case "healthy":
					state = "running"
				case "critical":
					state = "stopped"
				}
			}

			info := ServiceUpgradeInfo{
				Name:             name,
				DisplayName:      name,
				Category:         categorize(base),
				Kind:             kindForService(base),
				InstalledVersion: installed,
				State:            state,
				Port:             port,
				DerivedStatus:    derivedStatus,
			}

			latest, hasLatest := latestByName[base]
			if hasLatest {
				info.LatestVersion = latest.version
				info.LatestBuildNum = latest.buildNumber
			}

			// Determine upgrade availability.
			if installed == "" || !hasLatest {
				summary.Unknown++
			} else if cmp, err := versionutil.CompareFull(installed, 0, latest.version, latest.buildNumber); err == nil && cmp < 0 {
				info.UpdateAvailable = true
				summary.UpdatesAvailable++
			} else {
				summary.UpToDate++
			}

			services = append(services, info)
		}

		summary.TotalInstalled = len(services)
		for _, s := range services {
			switch s.Kind {
			case "INFRASTRUCTURE":
				summary.InfrastructureCount++
			case "APPLICATION":
				summary.ApplicationCount++
			default:
				summary.ServiceCount++
			}
		}

		// Determine repository status.
		repoStatus := "ok"
		if bndErr != nil {
			repoStatus = "unreachable"
		} else if len(bundles) == 0 {
			repoStatus = "empty"
		}

		resp := UpgradesStatusResponse{
			NowUnixMs:        time.Now().UnixMilli(),
			Node:             admin.Hostname(),
			Platform:         platform,
			Services:         services,
			Summary:          summary,
			RepositoryStatus: repoStatus,
		}

		if bndErr != nil {
			// Still return partial data — frontend can show "repository unreachable".
			resp.Summary.Unknown = resp.Summary.TotalInstalled
			resp.Summary.UpdatesAvailable = 0
			resp.Summary.UpToDate = 0
			for i := range resp.Services {
				resp.Services[i].LatestVersion = ""
				resp.Services[i].UpdateAvailable = false
			}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

type latestInfo struct {
	version     string
	buildNumber int64
}

// buildLatestIndex groups bundles by normalized service name and picks the
// highest version (semver + build number) for each, filtering to the given platform.
func buildLatestIndex(bundles []BundleInfo, platform string) map[string]latestInfo {
	result := make(map[string]latestInfo)
	for _, b := range bundles {
		if b.Platform != "" && b.Platform != platform {
			continue
		}
		base := normalizeServiceName(b.Name)
		if base == "" {
			continue
		}
		existing, ok := result[base]
		if !ok {
			result[base] = latestInfo{version: b.Version, buildNumber: b.BuildNumber}
			continue
		}
		cmp, err := versionutil.CompareFull(existing.version, existing.buildNumber, b.Version, b.BuildNumber)
		if err == nil && cmp < 0 {
			result[base] = latestInfo{version: b.Version, buildNumber: b.BuildNumber}
		}
	}
	return result
}

// normalizeServiceName extracts the base service name from full etcd/config names.
// Examples: "repository.PackageRepository" -> "repository"
//
//	"gateway" -> "gateway"
//	"dns.DnsService" -> "dns"
func normalizeServiceName(name string) string {
	name = strings.TrimSuffix(name, ".service")
	if idx := strings.IndexByte(name, '.'); idx > 0 {
		name = name[:idx]
	}
	if idx := strings.IndexByte(name, ':'); idx > 0 {
		name = name[:idx]
	}
	return strings.ToLower(name)
}

// ── Plan / Apply / Job Status ────────────────────────────────────────────────

// NodeAgentProvider gives the handler access to the local node-agent.
type NodeAgentProvider interface {
	ApplyPlanV1(ctx context.Context, plan *planpb.NodePlan) (operationID string, err error)
	GetPlanStatus(ctx context.Context, operationID string) (*planpb.NodePlanStatus, error)
}

// ControllerUpgradePlanItem is the plan item returned by the controller.
type ControllerUpgradePlanItem struct {
	Service         string
	FromVersion     string
	FromBuildNumber int64
	ToVersion       string
	ToBuildNumber   int64
	PackageName     string
	SHA256          string
	RestartRequired bool
	Impacts         []string
}

// ControllerPlanResult holds the result of a controller PlanServiceUpgrades call.
type ControllerPlanResult struct {
	Items            []ControllerUpgradePlanItem
	RepositoryStatus string // ok | unreachable | empty
}

// ControllerApplyResult holds the result of a controller ApplyServiceUpgrades call.
type ControllerApplyResult struct {
	OK          bool
	OperationID string
	Message     string
}

// ControllerProvider delegates upgrade planning and execution to the cluster controller.
type ControllerProvider interface {
	PlanServiceUpgrades(ctx context.Context, services []string) (*ControllerPlanResult, error)
	ApplyServiceUpgrades(ctx context.Context, services []string) (*ControllerApplyResult, error)
}

// ── Plan types ──────────────────────────────────────────────────────────────

// UpgradePlanRequest is the POST body for /admin/upgrades/plan.
type UpgradePlanRequest struct {
	Services []string `json:"services"` // base service names to upgrade
}

// UpgradePlanResponse shows what will happen if upgrades are applied.
type UpgradePlanResponse struct {
	Plan []UpgradePlanItem `json:"plan"`
}

// UpgradePlanItem describes a single service upgrade in the plan.
type UpgradePlanItem struct {
	Service         string   `json:"service"`
	Kind            string   `json:"kind,omitempty"` // "SERVICE", "APPLICATION", "INFRASTRUCTURE"
	From            string   `json:"from"`
	FromBuildNumber int64    `json:"from_build_number,omitempty"`
	To              string   `json:"to"`
	ToBuildNumber   int64    `json:"to_build_number,omitempty"`
	Package         string   `json:"package"`
	RestartRequired bool     `json:"restart_required"`
	Impacts         []string `json:"impacts,omitempty"`
}

// UpgradeApplyRequest is the POST body for /admin/upgrades/apply.
type UpgradeApplyRequest struct {
	Services []string `json:"services"` // base service names to upgrade
}

// UpgradeApplyResponse is returned when an upgrade is started.
type UpgradeApplyResponse struct {
	OK          bool   `json:"ok"`
	OperationID string `json:"operation_id"`
	Message     string `json:"message"`
}

// UpgradeJobResponse describes the state of an upgrade operation.
type UpgradeJobResponse struct {
	OperationID string           `json:"operation_id"`
	Status      string           `json:"status"` // pending | running | success | failed
	Steps       []UpgradeJobStep `json:"steps"`
	Progress    int              `json:"progress"` // 0-100
	Error       string           `json:"error,omitempty"`
}

// UpgradeJobStep describes one step in the upgrade job.
type UpgradeJobStep struct {
	ID      string `json:"id"`
	State   string `json:"state"` // pending | running | ok | failed | skipped
	Message string `json:"message,omitempty"`
}

// ── Impact map ──────────────────────────────────────────────────────────────

// ── Plan handler ────────────────────────────────────────────────────────────

// NewUpgradePlanHandler returns a POST handler for /admin/upgrades/plan.
// Planning is delegated to the cluster controller — gateway is a thin adapter.
func NewUpgradePlanHandler(controller ControllerProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpgradePlanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeUpgradeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid JSON: " + err.Error()})
			return
		}
		if len(req.Services) == 0 {
			writeUpgradeJSON(w, http.StatusBadRequest, map[string]any{"error": "at least one service is required"})
			return
		}

		result, err := controller.PlanServiceUpgrades(r.Context(), req.Services)
		if err != nil {
			writeUpgradeJSON(w, http.StatusInternalServerError, map[string]any{"error": "plan upgrades: " + err.Error()})
			return
		}

		// Map controller result to HTTP response.
		var plan []UpgradePlanItem
		for _, item := range result.Items {
			plan = append(plan, UpgradePlanItem{
				Service:         item.Service,
				From:            item.FromVersion,
				FromBuildNumber: item.FromBuildNumber,
				To:              item.ToVersion,
				ToBuildNumber:   item.ToBuildNumber,
				Package:         item.PackageName,
				RestartRequired: item.RestartRequired,
				Impacts:         item.Impacts,
			})
		}

		writeUpgradeJSON(w, http.StatusOK, UpgradePlanResponse{Plan: plan})
	})
}

// ── Apply handler ───────────────────────────────────────────────────────────

// NewUpgradeApplyHandler returns a POST handler for /admin/upgrades/apply.
// Upgrade execution is delegated to the cluster controller — gateway is a thin adapter.
// If a JobStore is provided, it persists the job record for history tracking.
func NewUpgradeApplyHandler(controller ControllerProvider, store ...*JobStore) http.Handler {
	var jobStore *JobStore
	if len(store) > 0 {
		jobStore = store[0]
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpgradeApplyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeUpgradeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid JSON: " + err.Error()})
			return
		}
		if len(req.Services) == 0 {
			writeUpgradeJSON(w, http.StatusBadRequest, map[string]any{"error": "at least one service is required"})
			return
		}

		// Delegate to cluster controller.
		result, err := controller.ApplyServiceUpgrades(r.Context(), req.Services)
		if err != nil {
			writeUpgradeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "apply upgrades: " + err.Error(),
			})
			return
		}

		if !result.OK {
			writeUpgradeJSON(w, http.StatusUnprocessableEntity, UpgradeApplyResponse{
				OK:      false,
				Message: result.Message,
			})
			return
		}

		// Persist job record for history.
		if jobStore != nil {
			// Get the plan to record service details.
			planResult, _ := controller.PlanServiceUpgrades(r.Context(), req.Services)
			var jobServices []UpgradeJobService
			if planResult != nil {
				for _, item := range planResult.Items {
					jobServices = append(jobServices, UpgradeJobService{
						Name: item.Service,
						From: item.FromVersion,
						To:   item.ToVersion,
					})
				}
			}
			_ = jobStore.Save(UpgradeJobRecord{
				OperationID: result.OperationID,
				StartedAt:   time.Now().UnixMilli(),
				Status:      "running",
				Services:    jobServices,
				IssuedBy:    "admin-ui",
			})
		}

		writeUpgradeJSON(w, http.StatusOK, UpgradeApplyResponse{
			OK:          result.OK,
			OperationID: result.OperationID,
			Message:     result.Message,
		})
	})
}

// ── Job status handler ──────────────────────────────────────────────────────

// NewUpgradeJobStatusHandler returns a GET handler for /admin/upgrades/jobs.
// Query param: ?id=<operation_id>
func NewUpgradeJobStatusHandler(agent NodeAgentProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		opID := r.URL.Query().Get("id")
		if opID == "" {
			writeUpgradeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing 'id' query parameter"})
			return
		}

		status, err := agent.GetPlanStatus(r.Context(), opID)
		if err != nil {
			writeUpgradeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "get plan status: " + err.Error(),
			})
			return
		}

		// Map proto status to response.
		resp := UpgradeJobResponse{
			OperationID: opID,
			Status:      planStateToString(status.GetState()),
			Error:       status.GetErrorMessage(),
		}

		totalSteps := len(status.GetSteps())
		doneSteps := 0
		for _, s := range status.GetSteps() {
			step := UpgradeJobStep{
				ID:      s.GetId(),
				State:   stepStateToString(s.GetState()),
				Message: s.GetMessage(),
			}
			resp.Steps = append(resp.Steps, step)
			if s.GetState() == planpb.StepState_STEP_OK || s.GetState() == planpb.StepState_STEP_SKIPPED {
				doneSteps++
			}
		}
		if totalSteps > 0 {
			resp.Progress = doneSteps * 100 / totalSteps
		}

		writeUpgradeJSON(w, http.StatusOK, resp)
	})
}

func planStateToString(s planpb.PlanState) string {
	switch s {
	case planpb.PlanState_PLAN_PENDING:
		return "pending"
	case planpb.PlanState_PLAN_RUNNING:
		return "running"
	case planpb.PlanState_PLAN_SUCCEEDED:
		return "success"
	case planpb.PlanState_PLAN_FAILED:
		return "failed"
	case planpb.PlanState_PLAN_ROLLING_BACK:
		return "rolling_back"
	case planpb.PlanState_PLAN_ROLLED_BACK:
		return "rolled_back"
	default:
		return "unknown"
	}
}

func stepStateToString(s planpb.StepState) string {
	switch s {
	case planpb.StepState_STEP_PENDING:
		return "pending"
	case planpb.StepState_STEP_RUNNING:
		return "running"
	case planpb.StepState_STEP_OK:
		return "ok"
	case planpb.StepState_STEP_FAILED:
		return "failed"
	case planpb.StepState_STEP_SKIPPED:
		return "skipped"
	default:
		return "unknown"
	}
}

func writeUpgradeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
