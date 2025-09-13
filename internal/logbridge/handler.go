package logbridge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/log/logpb"
	Utility "github.com/globulario/utility"
	"google.golang.org/protobuf/encoding/protojson"
)

// Transport topic for log events
const TopicNewLog = "new_log_evt"

// EventSlogHandler converts slog records to logpb.LogInfo and publishes them
// over the Event service. It coalesces identical lines into a single message
// with Occurences>1 to avoid flooding the bus.
type EventSlogHandler struct {
	ec       *event_client.Event_Client
	app      string // LogInfo.application
	minLevel slog.Level
	mirror   slog.Handler // optional local mirror (stderr/file)

	// coalescing
	mu     sync.Mutex
	bucket map[string]*coalesced // key -> aggregation
	flush  *time.Ticker

	stop chan struct{}
}

type coalesced struct {
	info *logpb.LogInfo
	seen int64
}

// Opts configures the handler.
type Opts struct {
	App           string
	MinLevel      slog.Level
	Mirror        slog.Handler
	FlushInterval time.Duration // e.g. 500ms–2s
}

// New wires the handler to the Event service at address and starts the flusher.
func New(_ context.Context, address string, o Opts) (*EventSlogHandler, error) {
	Utility.RegisterFunction("NewEventService_Client", event_client.NewEventService_Client)
	c, err := globular_client.GetClient(address, "event.EventService", "NewEventService_Client")
	if err != nil {
		return nil, err
	}
	h := &EventSlogHandler{
		ec:       c.(*event_client.Event_Client),
		app:      o.App,
		mirror:   o.Mirror,
		minLevel: o.MinLevel,
		bucket:   make(map[string]*coalesced),
		flush:    time.NewTicker(max(o.FlushInterval, 750*time.Millisecond)),
		stop:     make(chan struct{}),
	}
	go h.flusher()
	return h, nil
}

// Close stops the flusher and publishes any pending coalesced entries.
func (h *EventSlogHandler) Close() {
	select {
	case <-h.stop:
		// already closed
	default:
		close(h.stop)
	}
	h.flush.Stop()

	// final flush
	h.mu.Lock()
	batch := h.bucket
	h.bucket = make(map[string]*coalesced)
	h.mu.Unlock()

	for _, c := range batch {
		c.info.Occurences = c.seen
		data, _ := protojson.Marshal(c.info)
		_ = h.ec.Publish(TopicNewLog, data)
	}
}

// Enabled implements slog.Handler.
func (h *EventSlogHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.minLevel
}

// Handle implements slog.Handler.
func (h *EventSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	// Mirror locally if provided
	if h.mirror != nil && h.mirror.Enabled(ctx, r.Level) {
		_ = h.mirror.Handle(ctx, r)
	}

	method, line := extractSource(r)         // method and file:line if AddSource=true
	component, fields := splitAttrs(nil, &r) // attrs from the Record
	id := deriveID(h.app, component, method, line, r.Message)

	// Prefer the record's timestamp; fall back to now
	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}

	info := &logpb.LogInfo{
		Id:          id,
		Level:       toLogLevel(r.Level),
		Application: h.app,
		Method:      method,
		Line:        line,
		Message:     r.Message,
		Occurences:  1,
		TimestampMs: ts.UnixMilli(),
		Component:   component,
	}
	if len(fields) > 0 {
		info.Fields = fields
	}

	// Aggregate short bursts of identical logs into one message (Occurences++)
	h.mu.Lock()
	if agg, ok := h.bucket[id]; ok && samePayload(agg.info, info) {
		agg.seen++
	} else {
		h.bucket[id] = &coalesced{info: info, seen: 1}
	}
	h.mu.Unlock()
	return nil
}

// WithAttrs returns a lightweight wrapper that attaches the given attrs.
func (h *EventSlogHandler) WithAttrs(as []slog.Attr) slog.Handler {
	var m slog.Handler
	if h.mirror != nil {
		m = h.mirror.WithAttrs(as)
	}
	return &attached{
		parent: h,
		attrs:  as,
		mirror: m,
	}
}

// WithGroup returns a lightweight wrapper that records the group.
func (h *EventSlogHandler) WithGroup(name string) slog.Handler {
	var m slog.Handler
	if h.mirror != nil {
		m = h.mirror.WithGroup(name)
	}
	return &attached{
		parent: h,
		groups: []string{name},
		mirror: m,
	}
}

// attached is a thin wrapper that carries extra attrs/groups and delegates to the parent.
type attached struct {
	parent *EventSlogHandler
	attrs  []slog.Attr
	groups []string
	mirror slog.Handler
}

func (a *attached) Enabled(ctx context.Context, lvl slog.Level) bool {
	return a.parent.Enabled(ctx, lvl)
}

func (a *attached) Handle(ctx context.Context, r slog.Record) error {
	// Mirror first if present
	if a.mirror != nil && a.mirror.Enabled(ctx, r.Level) {
		_ = a.mirror.Handle(ctx, r)
	} else if a.parent.mirror != nil && a.parent.mirror.Enabled(ctx, r.Level) {
		_ = a.parent.mirror.Handle(ctx, r)
	}

	// Combine attached attrs with record attrs
	component, fields := splitAttrs(a.attrs, &r)

	// Reuse parent's coalescing & publish path
	method, line := extractSource(r)
	id := deriveID(a.parent.app, component, method, line, r.Message)

	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}

	info := &logpb.LogInfo{
		Id:          id,
		Level:       toLogLevel(r.Level),
		Application: a.parent.app,
		Method:      method,
		Line:        line,
		Message:     r.Message,
		Occurences:  1,
		TimestampMs: ts.UnixMilli(),
		Component:   component,
	}
	if len(fields) > 0 {
		info.Fields = fields
	}

	a.parent.mu.Lock()
	if agg, ok := a.parent.bucket[id]; ok && samePayload(agg.info, info) {
		agg.seen++
	} else {
		a.parent.bucket[id] = &coalesced{info: info, seen: 1}
	}
	a.parent.mu.Unlock()
	return nil
}

func (a *attached) WithAttrs(as []slog.Attr) slog.Handler {
	var m slog.Handler
	if a.mirror != nil {
		m = a.mirror.WithAttrs(as)
	}
	na := &attached{
		parent: a.parent,
		mirror: m,
	}
	na.attrs = append(append([]slog.Attr{}, a.attrs...), as...)
	na.groups = append([]string{}, a.groups...)
	return na
}

func (a *attached) WithGroup(name string) slog.Handler {
	var m slog.Handler
	if a.mirror != nil {
		m = a.mirror.WithGroup(name)
	}
	na := &attached{
		parent: a.parent,
		mirror: m,
	}
	na.attrs = append([]slog.Attr{}, a.attrs...)
	na.groups = append(append([]string{}, a.groups...), name)
	return na
}

// Periodically publish coalesced entries.
func (h *EventSlogHandler) flusher() {
	for {
		select {
		case <-h.stop:
			return
		case <-h.flush.C:
			h.mu.Lock()
			batch := h.bucket
			h.bucket = make(map[string]*coalesced)
			h.mu.Unlock()

			for _, c := range batch {
				c.info.Occurences = c.seen
				data, _ := protojson.Marshal(c.info)
				_ = h.ec.Publish(TopicNewLog, data) // best-effort; don’t block request paths
			}
		}
	}
}

func toLogLevel(l slog.Level) logpb.LogLevel {
	switch {
	case l >= slog.LevelError:
		return logpb.LogLevel_ERROR_MESSAGE
	case l >= slog.LevelWarn:
		return logpb.LogLevel_WARN_MESSAGE
	case l >= slog.LevelInfo:
		return logpb.LogLevel_INFO_MESSAGE
	default:
		return logpb.LogLevel_DEBUG_MESSAGE
	}
}

func extractSource(r slog.Record) (method, line string) {
	// Works when slog.HandlerOptions{AddSource:true} is set on the logger/handler.
	if r.PC == 0 {
		return "", ""
	}
	fn := runtime.FuncForPC(r.PC)
	if fn == nil {
		return "", ""
	}
	file, ln := fn.FileLine(r.PC)
	name := fn.Name()
	// Shorten func name (strip package path)
	if i := strings.LastIndexFunc(name, func(r rune) bool { return r == '/' || r == '.' }); i >= 0 && i+1 < len(name) {
		name = name[i+1:]
	}
	return name, fmt.Sprintf("%s:%d", shortFile(file), ln)
}

func shortFile(f string) string {
	if i := strings.LastIndexAny(f, "/\\"); i >= 0 && i+1 < len(f) {
		return f[i+1:]
	}
	return f
}

// splitAttrs pulls component and fields from attached + record attrs.
// Recognizes "component" and "comp" as the component tag.
func splitAttrs(attached []slog.Attr, r *slog.Record) (component string, fields map[string]string) {
	fields = make(map[string]string)

	apply := func(a slog.Attr) {
		// Manually resolve lazy values if present
		if a.Value.Kind() == slog.KindLogValuer {
			if lv, ok := a.Value.Any().(slog.LogValuer); ok {
				a.Value = lv.LogValue()
			}
		}
		switch a.Key {
		case "component", "comp":
			component = attrToString(a.Value)
		default:
			fields[a.Key] = attrToString(a.Value)
		}
	}

	for _, a := range attached {
		apply(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		apply(a)
		return true
	})

	// Remove empty map to keep payload small
	if len(fields) == 0 {
		fields = nil
	}
	return component, fields
}

func attrToString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339Nano)
	default:
		return fmt.Sprint(v.Any())
	}
}

func deriveID(app, comp, method, line, msg string) string {
	h := sha256.Sum256([]byte(app + "|" + comp + "|" + method + "|" + line + "|" + msg))
	return hex.EncodeToString(h[:16]) // short, stable id for same message at same site
}

func samePayload(a, b *logpb.LogInfo) bool {
	return a.Application == b.Application &&
		a.Component == b.Component &&
		a.Method == b.Method &&
		a.Line == b.Line &&
		a.Message == b.Message &&
		a.Level == b.Level
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
