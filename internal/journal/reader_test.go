package journal

import (
	"reflect"
	"testing"
)

func TestSplitLinesReturnsTailAndDropsEmptyLines(t *testing.T) {
	input := `one

 two
three
four
`
	got := splitLines(input, 2)
	want := []string{"three", "four"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitLines=%v want %v", got, want)
	}
}
