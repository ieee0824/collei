package aggregator

import (
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"
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

func TestContainer_add(t *testing.T) {
	tests := []struct {
		name string
		opt  *containerOpt
		in   []string
		want int
	}{
		{
			name: "毎回送信される",
			opt: &containerOpt{
				emitDuration: 1000 * time.Hour,
				maxCount:     0,
			},
			in: lo.Map([]int{0, 1, 2, 3}, func(i int, _ int) string {
				return fmt.Sprint(i)
			}),
		}, {
			name: "カウントを超えないので送信されない",
			opt: &containerOpt{
				emitDuration: 1000 * time.Hour,
				maxCount:     1000,
			},
			in: lo.Map([]int{0, 1, 2, 3}, func(i int, _ int) string {
				return fmt.Sprint(i)
			}),
			want: 4,
		}, {
			name: "一部送信される",
			opt: &containerOpt{
				emitDuration: 1000 * time.Hour,
				maxCount:     3,
			},
			in: lo.Map([]int{0, 1, 2, 3}, func(i int, _ int) string {
				return fmt.Sprint(i)
			}),
			want: 1,
		}, {
			name: "全部送信された",
			opt: &containerOpt{
				emitDuration: 1000 * time.Hour,
				maxCount:     4,
			},
			in: lo.Map([]int{0, 1, 2, 3}, func(i int, _ int) string {
				return fmt.Sprint(i)
			}),
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testContainer := newContainer[string](test.opt)

			for _, v := range test.in {
				testContainer.add(v)
			}

			assert.Equal(t, test.want, len(testContainer.elems))
		})
	}
}
