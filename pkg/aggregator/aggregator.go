package aggregator

import (
	"io"
	"time"
)

type ParserFunc[T any] func([]byte) (*T, error)
type KeyGeneratorFunc[T any] func(*T) (string, error)

type Aggregator[T any] interface {
	io.Writer
}

func New[T any](w io.Writer, opt ...OptFunc[T]) Aggregator[T] {
	op := &Opt[T]{
		Format:       "json",
		MaxCunt:      10,
		EmitDuration: 1 * time.Minute,
	}
	for _, f := range opt {
		f(op)
	}
	return newWriter(w, op)
}
