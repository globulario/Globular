package logsink

import (
	"bytes"
	"strings"
	"testing"

	"github.com/globulario/services/golang/log/logpb"
)

func TestSplitFlushablePreservesTrailingFragment(t *testing.T) {
	var keep string
	got := splitFlushable("first\nsecond\rthird", &keep)
	want := []string{"first", "second"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("splitFlushable=%v want %v", got, want)
	}
	if keep != "third" {
		t.Fatalf("keep=%q want third", keep)
	}
}

func TestServiceLogWriterEchoesCompletedLinesWithoutRemoteClient(t *testing.T) {
	var echo bytes.Buffer
	w := NewServiceLogWriter("127.0.0.1:1", "app", "user", "method", logpb.LogLevel_INFO_MESSAGE, &echo)

	if n, err := w.Write([]byte("hello")); err != nil || n != 5 {
		t.Fatalf("first Write n=%d err=%v", n, err)
	}
	if echo.Len() != 0 {
		t.Fatalf("echo should stay empty for partial line, got %q", echo.String())
	}

	payload := " world\nsecond\rthird"
	if n, err := w.Write([]byte(payload)); err != nil || n != len(payload) {
		t.Fatalf("second Write n=%d err=%v", n, err)
	}
	if got := echo.String(); got != "hello world\nsecond\n" {
		t.Fatalf("echo after second write=%q", got)
	}

	payload = " tail\n"
	if n, err := w.Write([]byte(payload)); err != nil || n != len(payload) {
		t.Fatalf("third Write n=%d err=%v", n, err)
	}
	if got := echo.String(); got != "hello world\nsecond\nthird tail\n" {
		t.Fatalf("final echo=%q", got)
	}
}
