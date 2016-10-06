package trace

import (
	"bytes"
	"testing"
)

// Testで始まり *testing.T を１つ受け取る関数はユニットテストとみなされる。
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s' という誤った文字列が出力されました", buf.String())
		}
	}
}
