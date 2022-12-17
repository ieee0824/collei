package aggregator

import (
	"fmt"
	"io"
	"time"

	"github.com/samber/lo"
)

type ParserFunc[T any] func([]byte) (*T, error)
type KeyGeneratorFunc[T any] func(*T) (string, error)

type Aggregator[T any] interface {
	io.Writer
}

func New[T any](w io.Writer, optFs ...OptFunc[T]) Aggregator[T] {
	// set default
	op := &Opt[T]{
		Format:       "json",
		MaxCnt:       10,
		EmitDuration: 1 * time.Minute,
		KeyGenerator: func(t *T) (string, error) {
			return fmt.Sprint(*t), nil
		},
	}
	lo.ForEach(optFs, func(optF OptFunc[T], _ int) {
		optF(op)
	})
	return newWriter(w, op)
}
