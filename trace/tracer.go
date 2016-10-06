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

type tracer struct {
	out io.Writer
}

// Writeメソッドを実装していれば何でも渡せる
// インタフェース「Tracer」（の実装オブジェクト）を生成してリターンする関数
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// 構造体「tracer」をレシーバとして、インタフェース「Tracer」を実装
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
