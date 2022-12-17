package aggregator

import (
	"encoding/json"
	"io"
	"log"
	"time"
)

func jsonParser[T any](b []byte) (*T, error) {
	var ret T

	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

type Opt[T any] struct {
	Format       string
	KeyGenerator KeyGeneratorFunc[T]
	MaxCnt       int
	EmitDuration time.Duration
}

type OptFunc[T any] func(o *Opt[T])

func newWriter[T any](w io.Writer, opt *Opt[T]) *writer[T] {
	rw := &writer[T]{
		containers:           map[string]*container[T]{},
		maxCount:             opt.MaxCnt,
		emitDuration:         opt.EmitDuration,
		aggregatedLogStrPipe: make(chan string, 65536),
		w:                    w,
	}

	switch opt.Format {
	case "json":
		rw.parse = jsonParser[T]
	}
	rw.keyGenerate = opt.KeyGenerator

	go rw.emit()

	return rw
}

type writer[T any] struct {
	w                    io.Writer
	containers           map[string]*container[T]
	parse                ParserFunc[T]
	keyGenerate          KeyGeneratorFunc[T]
	maxCount             int
	emitDuration         time.Duration
	aggregatedLogStrPipe chan string
}

func (impl *writer[T]) emit() {
	for {
		str, ok := <-impl.aggregatedLogStrPipe
		if ok {
			_, err := impl.w.Write([]byte(str))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (impl *writer[T]) Write(b []byte) (int, error) {
	raw, err := impl.parse(b)
	if err != nil {
		return 0, err
	}
	key, err := impl.keyGenerate(raw)
	if err != nil {
		return 0, err
	}
	_, exists := impl.containers[key]
	if !exists {
		impl.containers[key] = newContainer[T](
			key,
			impl.aggregatedLogStrPipe,
			&containerOpt{
				maxCount:     impl.maxCount,
				emitDuration: impl.emitDuration,
			},
		)
	}
	impl.containers[key].add(raw)

	return len(b), nil
}
