package aggregator

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func newContainer[T any](pipe chan<- string, opt *containerOpt) *container[T] {
	return &container[T]{
		mut:                  sync.Mutex{},
		elems:                []*T{},
		aggregatedLogStrPipe: pipe,
		maxCount:             opt.getMaxCount(),
		emitDuration:         opt.getEmitDuration(),
	}
}

type container[T any] struct {
	mut                  sync.Mutex
	elems                []*T
	aggregatedLogStrPipe chan<- string
	maxCount             int
	emitDuration         time.Duration
}

func (impl *container[T]) add(t T) {
	impl.mut.Lock()
	defer impl.mut.Unlock()
	impl.elems = append(impl.elems, &t)
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

func (impl *container[T]) string() string {
	if len(impl.elems) == 0 {
		return ""
	}

	strs := lo.Map(impl.elems, func(e *T, _ int) string {
		if e == nil {
			return ""
		}

		return fmt.Sprintf("%v", *e)
	})

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(lo.Filter(strs, func(s string, _ int) bool {
		return s != ""
	}))

	return buf.String()
}
