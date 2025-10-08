// package logsink (or a new tiny package in Globular)
// This avoids importing it from Services; Services just takes io.Writer.

package logsink

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/log/log_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
)

type ServiceLogWriter struct {
	mu      sync.Mutex
	buf     bytes.Buffer
	app     string
	user    string
	method  string
	level   logpb.LogLevel
	echoTo  io.Writer // optional: also print to local console
	lc      *log_client.Log_Client
	address string
}

// NewServiceLogWriter returns a writer that sends each flushed line to LogService at `address`.
func NewServiceLogWriter(address, application, user, method string, level logpb.LogLevel, echo io.Writer) *ServiceLogWriter {
	return &ServiceLogWriter{
		app:     application,
		user:    user,
		method:  method,
		level:   level,
		echoTo:  echo,
		address: address,
	}
}

func (w *ServiceLogWriter) ensureClient() error {
	if w.lc != nil {
		return nil
	}
	Utility.RegisterFunction("NewLogService_Client", log_client.NewLogService_Client)
	c, err := globular_client.GetClient(w.address, "log.LogService", "NewLogService_Client")
	if err != nil {
		return err
	}
	w.lc = c.(*log_client.Log_Client)
	return nil
}

func (w *ServiceLogWriter) Write(p []byte) (int, error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	// Append and flush on newline or CR
	n, _ := w.buf.Write(p)
	data := w.buf.Bytes()

	// Split on both \n and \r (so CR-only progress lines still flush)
	// We’ll keep trailing partial line in buffer.
	s := string(data)
	var keep string
	lines := splitFlushable(s, &keep)
	w.buf.Reset()
	if keep != "" {
		w.buf.WriteString(keep)
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		// Optionally echo to console
		if w.echoTo != nil {
			_, _ = w.echoTo.Write([]byte(line + "\n"))
		}

		// Get the local token
		mac, _ := config.GetMacAddress()
		localToken, _ := security.GetLocalToken(mac)

		// Send to LogService (best-effort)
		if err := w.ensureClient(); err == nil {
			_ = w.lc.Log(w.app, w.user, w.method, w.level, line, "-1", "-1", localToken)
		}
	}
	return n, nil
}

func splitFlushable(s string, keep *string) []string {
	// Flush on \n or \r; keep any trailing fragment without either
	// Examples: "abc\n" => ["abc"], keep=""
	//           "abc\rdef" => ["abc"], keep="def"
	//           "abc" => [], keep="abc"
	last := -1
	out := []string{}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			out = append(out, s[last+1:i])
			last = i
		}
	}
	if last+1 < len(s) {
		*keep = s[last+1:]
	} else {
		*keep = ""
	}
	return out
}
