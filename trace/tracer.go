package trace

import (
	"fmt"
	"io"
)

// コード内での出来事を記録できるオブジェクトを表すインタフェース
// 先頭が大文字のため、公開される
type Tracer interface {
	Trace(...interface{}) // 任意の型の引数を何個でも（ゼロ個でも）受け取るメソッド
}

// 小文字始まりのため公開されない構造体
type tracer struct {
	out io.Writer
}

type nilTracer struct{}

// Writeメソッドを実装していれば何でも渡せる
// インタフェース「Tracer」（に合致したオブジェクト）を生成してリターンする関数
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// 構造体「tracer」をレシーバとして、インタフェース「Tracer」を実装
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

// 構造体「nilTracer」をレシーバとして、インタフェース「Tracer」を実装
func (t *nilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}
