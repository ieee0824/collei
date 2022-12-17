package aggregator

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/samber/lo"
)

type containerOpt struct {
	maxCount     int
	emitDuration time.Duration
}

func (impl *containerOpt) getMaxCount() int {
	if impl == nil {
		return 10
	}
	return impl.maxCount
}

func (impl *containerOpt) getEmitDuration() time.Duration {
	if impl == nil {
		return 1 * time.Minute
	}
	return impl.emitDuration
}

func newContainer[T any](key string, pipe chan<- string, opt *containerOpt) *container[T] {
	c := &container[T]{
		key:                  key,
		mut:                  sync.Mutex{},
		elems:                []*T{},
		aggregatedLogStrPipe: pipe,
		maxCount:             opt.getMaxCount(),
		emitDuration:         opt.getEmitDuration(),
	}
	go c.emitUseDuration()

	return c
}

type container[T any] struct {
	mut                  sync.Mutex
	elems                []*T
	aggregatedLogStrPipe chan<- string
	maxCount             int
	emitDuration         time.Duration
	key                  string
}

func (impl *container[T]) add(t *T) {
	impl.mut.Lock()
	defer impl.mut.Unlock()
	impl.elems = append(impl.elems, t)
	if len(impl.elems) < impl.maxCount {
		return
	}
	impl.aggregatedLogStrPipe <- impl.string()
	impl.elems = []*T{}
}

func (impl *container[T]) emitUseDuration() {
	ticker := time.NewTicker(impl.emitDuration)
	for {
		<-ticker.C
		func() {
			impl.mut.Lock()
			defer impl.mut.Unlock()
			if len(impl.elems) == 0 {
				return
			}
			impl.aggregatedLogStrPipe <- impl.string()
			impl.elems = []*T{}
		}()
	}
}

type output struct {
	Key  string `json:"key"`
	Logs any    `json:"logs"`
}

func (impl *container[T]) string() string {
	if len(impl.elems) == 0 {
		return ""
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(&output{
		Key: impl.key,
		Logs: lo.Filter(impl.elems, func(e *T, _ int) bool {
			return e != nil
		}),
	})

	return buf.String()
}
