package logsink

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/log/log_client"
	"github.com/globulario/services/golang/log/logpb"
	Utility "github.com/globulario/utility"
	"github.com/gookit/color"
	"google.golang.org/protobuf/encoding/protojson"
)

// Filter controls which logs are displayed and how.
type Filter struct {
	MinLevel   logpb.LogLevel
	Apps       map[string]bool // nil => all apps (see BackfillApps note)
	Components map[string]bool // nil => all components
	ShowFields bool            // print structured fields after the message

	// Optional "catch up" settings: if BackfillSince > 0, the sink will fetch
	// and print recent logs once right after subscription is established.
	BackfillSince       time.Duration // e.g. 2*time.Minute
	BackfillPerAppLimit int           // e.g. 250
	BackfillApps        []string      // if empty and Apps==nil, a default list of known services is used
}

// ConsoleSink subscribes to new_log_evt and renders logs to the console.
type ConsoleSink struct {
	address string
	ec      *event_client.Event_Client
	filter  Filter
	stop    context.CancelFunc
	subUUID string

	backfilled bool // ensure we only backfill once per process
}

// NewConsoleSink creates a sink. It does not connect until Start() is called.
func NewConsoleSink(address string, filter Filter) *ConsoleSink {
	return &ConsoleSink{
		address: address,
		filter:  filter,
		// subUUID must be stable for the life of the process, but unique across processes
		subUUID: "globule-console-" + Utility.RandomUUID(),
	}
}

// Start launches a background retry loop that subscribes to EventService once
// it becomes available. The returned stop func cleanly unsubscribes.
func (s *ConsoleSink) Start() (func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	s.stop = cancel
	go s.run(ctx) // retry loop; see run()
	return func() { cancel() }, nil
}

func (s *ConsoleSink) run(ctx context.Context) {
	const topic = "new_log_evt"
	backoff := time.Second

	for {
		if ctx.Err() != nil {
			return
		}

		// Ensure we have an Event client
		if s.ec == nil {
			Utility.RegisterFunction("NewEventService_Client", event_client.NewEventService_Client)
			c, err := globular_client.GetClient(s.address, "event.EventService", "NewEventService_Client")
			if err != nil {
				// Event service not up yet: backoff & retry
				if !sleepOrDone(ctx, backoff) {
					return
				}
				backoff = grow(backoff, 5*time.Second)
				continue
			}
			s.ec = c.(*event_client.Event_Client)
		}

		// Try to subscribe (event_client will handle stream reconnects once subscribed)
		if err := s.ec.Subscribe(topic, s.subUUID, s.onEvent); err != nil {
			// Could be a race (dial, transient Unavailable). Reset client and retry.
			s.ec = nil
			if !sleepOrDone(ctx, backoff) {
				return
			}
			backoff = grow(backoff, 5*time.Second)
			continue
		}

		// We’re subscribed: kick a one-shot backfill if requested.
		if s.filter.BackfillSince > 0 && !s.backfilled {
			s.backfilled = true
			go s.backfillOnce()
		}

		// Subscribed: block until we’re asked to stop.
		<-ctx.Done()
		// Best-effort unsubscribe
		if s.ec != nil {
			_ = s.ec.UnSubscribe(topic, s.subUUID)
		}
		return
	}
}

func (s *ConsoleSink) onEvent(evt *eventpb.Event) {
	info := new(logpb.LogInfo)
	opts := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := opts.Unmarshal(evt.Data, info); err != nil {
		return
	}

	// Filtering
	if info.Level < s.filter.MinLevel {
		return
	}
	if s.filter.Apps != nil && !s.filter.Apps[info.Application] {
		return
	}
	if s.filter.Components != nil && !s.filter.Components[info.Component] {
		return
	}

	s.render(info)
}

// ---------- backfill (optional) ----------

func (s *ConsoleSink) backfillOnce() {
	// 1) Resolve which apps to query.
	apps := s.filter.BackfillApps
	if len(apps) == 0 && s.filter.Apps != nil {
		for a, ok := range s.filter.Apps {
			if ok {
				apps = append(apps, a)
			}
		}
		sort.Strings(apps)
	}
	if len(apps) == 0 {
		// Sensible defaults
		apps = []string{
			"event.EventService",
			"authentication.AuthenticationService",
			"log.LogService",
			"persistence.PersistenceService",
			"resource.ResourceService",
			"rbac.RbacService",
			"admin.AdminService",
			"discovery.PackageDiscovery",
			"applications_manager.ApplicationManagerService",
			"blog.BlogService",
			"catalog.CatalogService",
			"conversation.ConversationService",
			"dns.DnsService",
			"echo.EchoService",
			"file.FileService",
			"ldap.LdapService",
			"media.MediaService",
			"monitoring.MonitoringService",
			"repository.PackageRepository",
			"search.SearchService",
			"storage.StorageService",
			"title.TitleService",
			"torrent.TorrentService",
		}
	}

	// 2) Levels to query based on MinLevel (inclusive).
	levelNames := make([]string, 0, 4)
	for _, l := range []logpb.LogLevel{
		logpb.LogLevel_ERROR_MESSAGE,
		logpb.LogLevel_WARN_MESSAGE,
		logpb.LogLevel_INFO_MESSAGE,
		logpb.LogLevel_DEBUG_MESSAGE,
	} {
		if l >= s.filter.MinLevel {
			levelNames = append(levelNames, levelNameForQuery(l))
		}
	}
	if len(levelNames) == 0 {
		// Nothing to query at/above MinLevel.
		return
	}

	// 3) Get a LogService client (with short retry while it starts).
	Utility.RegisterFunction("NewLogService_Client", log_client.NewLogService_Client)

	deadline := time.Now().Add(20 * time.Second)
	var lc *log_client.Log_Client
	for {
		c, err := globular_client.GetClient(s.address, "log.LogService", "NewLogService_Client")
		if err == nil {
			lc = c.(*log_client.Log_Client)
			break
		}
		if time.Now().After(deadline) {
			color.Gray.Println("console backfill: LogService not ready, skipping")
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 4) Window + limit
	sinceMs := time.Now().Add(-s.filter.BackfillSince).UnixMilli()
	limit := s.filter.BackfillPerAppLimit
	if limit <= 0 {
		limit = 250
	}

	// 5) Query each (level, app) once.
	for _, app := range apps {
		for _, lvl := range levelNames {
			q := fmt.Sprintf("/%s/%s/*?since=%d&order=asc&limit=%d", lvl, app, sinceMs, limit)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			infos, err := lc.GetLogCtx(ctx, q)
			cancel() // ensure we don’t leak the context

			if err != nil {
				// Real error only.
				fmt.Printf("console backfill: query failed %q: %v\n", q, err)
				continue
			}

			// Respect component filter during backfill as well.
			filtered := 0
			for _, li := range infos {
				if s.filter.Components != nil && !s.filter.Components[li.Component] {
					continue
				}
				filtered++
				s.render(li)
			}

			// Optional: small, non-noisy summary when nothing matched.
			if len(infos) == 0 || filtered == 0 {
				color.Gray.Printf("console backfill: 0 entries for %q\n", q)
			}
		}
	}
}

// ---------- shared rendering ----------

func (s *ConsoleSink) render(info *logpb.LogInfo) {
	// Timestamp (prefer producer)
	ts := time.Now()
	if info.TimestampMs != 0 {
		ts = time.UnixMilli(info.TimestampMs)
	}

	// Trim full path to basename if present
	line := info.Line
	if line != "" {
		if i := strings.LastIndexAny(line, "/\\"); i >= 0 && i+1 < len(line) {
			// Keep "file.go:NNN"
			if j := strings.LastIndex(line, ":"); j > i+1 {
				line = line[i+1:]
			} else {
				line = line[i+1:]
			}
		}
	} else {
		line = "-"
	}

	comp := info.Component
	if comp == "" {
		comp = "-"
	}

	// Include component as its own column (kept narrow)
	header := fmt.Sprintf("[%s] %-24s | %-5s | %-7s | %s:%s",
		ts.Format("15:04:05.000"),
		clip(info.Application, 24),
		level(info.Level),
		clip(comp, 7),
		nz(info.Method), nz(line),
	)

	switch info.Level {
	case logpb.LogLevel_ERROR_MESSAGE:
		color.Red.Println(header)
	case logpb.LogLevel_WARN_MESSAGE:
		color.Yellow.Println(header)
	case logpb.LogLevel_INFO_MESSAGE:
		color.Cyan.Println(header)
	default:
		color.Gray.Println(header)
	}

	if msg := buildMsg(info.Message, int(info.Occurences)); msg != "" {
		color.Gray.Println(msg)
	}
	if s.filter.ShowFields && len(info.Fields) > 0 {
		color.Gray.Println(formatFields(info.Fields))
	}
}

// ---------- small helpers ----------

func buildMsg(s string, n int) string {
	if s == "" {
		return ""
	}
	if n > 1 {
		return fmt.Sprintf("%s  (x%d)", s, n)
	}
	return s
}

func level(l logpb.LogLevel) string {
	switch l {
	case logpb.LogLevel_ERROR_MESSAGE:
		return "ERROR"
	case logpb.LogLevel_WARN_MESSAGE:
		return "WARN"
	case logpb.LogLevel_INFO_MESSAGE:
		return "INFO"
	default:
		return "DEBUG"
	}
}

func levelNameForQuery(l logpb.LogLevel) string {
	switch l {
	case logpb.LogLevel_ERROR_MESSAGE:
		return "error"
	case logpb.LogLevel_WARN_MESSAGE:
		return "warning"
	case logpb.LogLevel_INFO_MESSAGE:
		return "info"
	case logpb.LogLevel_DEBUG_MESSAGE:
		return "debug"
	case logpb.LogLevel_TRACE_MESSAGE:
		return "trace"
	case logpb.LogLevel_FATAL_MESSAGE:
		return "fatal"
	default:
		return "info"
	}
}

func nz(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func clip(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func formatFields(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := "fields:"
	for _, k := range keys {
		out += " " + k + "=" + m[k]
	}
	return out
}

func sleepOrDone(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func grow(cur, max time.Duration) time.Duration {
	cur *= 2
	if cur > max {
		return max
	}
	return cur
}
