package aggregator

import (
	"io"
)

type ParserFunc[T any] func([]byte) (T, error)
type KeyGeneratorFunc[T any] func(T) (string, error)

type Aggregator[T any] interface {
	io.Writer
}

func New[T any]() Aggregator[T] {
	return newWriter[T](make(chan string, 1024))
}
