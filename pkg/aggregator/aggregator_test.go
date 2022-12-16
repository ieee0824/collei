package aggregator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContainerOpt_getMaxCount(t *testing.T) {
	tests := []struct {
		name string
		opt  *containerOpt
		want int
	}{
		{
			name: "optがnil",
			want: 10,
		},
		{
			name: "optが指定されている",
			opt: &containerOpt{
				maxCount: 100,
			},
			want: 100,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.opt.getMaxCount()
			assert.Equal(t, test.want, result)
		})
	}
}

func TestContainer_getEmitDuration(t *testing.T) {
	tests := []struct {
		name string
		opt  *containerOpt
		want time.Duration
	}{
		{
			name: "optがnil",
			want: 1 * time.Minute,
		},
		{
			name: "optが指定されている",
			want: 30 * time.Second,
			opt: &containerOpt{
				emitDuration: 30 * time.Second,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.opt.getEmitDuration()
			assert.Equal(t, test.want, result)
		})
	}
}

func TestNewContainer(t *testing.T) {
	tests := []struct {
		name string
		opt  *containerOpt
		want *container[any]
	}{
		{
			name: "optがnil",
			want: &container[any]{
				maxCount:     10,
				emitDuration: 1 * time.Minute,
			},
		},
		{
			name: "optが指定されている",
			opt: &containerOpt{
				maxCount:     100,
				emitDuration: 100 * time.Minute,
			},
			want: &container[any]{
				maxCount:     100,
				emitDuration: 100 * time.Minute,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := newContainer[any](test.opt)

			assert.NotNil(t, result.elems)
			assert.Equal(t, 0, len(result.elems))
			assert.NotNil(t, result.aggregatedLogStrPipe)
			assert.Equal(t, test.want.maxCount, result.maxCount)
			assert.Equal(t, test.want.emitDuration, result.emitDuration)
		})
	}
}
