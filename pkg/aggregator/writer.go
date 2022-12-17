package aggregator

import "time"

func newWriter[T any](pipe chan<- string) *writer[T] {
	return &writer[T]{
		containers:           map[string]*container[T]{},
		maxCount:             10,
		emitDuration:         1 * time.Minute,
		aggregatedLogStrPipe: pipe,
	}
}

type writer[T any] struct {
	containers           map[string]*container[T]
	parse                ParserFunc[T]
	keyGenerate          KeyGeneratorFunc[T]
	maxCount             int
	emitDuration         time.Duration
	aggregatedLogStrPipe chan<- string
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
