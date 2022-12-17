package aggregator

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opt  OptFunc[string]
		want *writer[string]
	}{
		{
			name: "opt funcを指定しない",
			want: &writer[string]{
				maxCount:     10,
				emitDuration: 1 * time.Minute,
			},
		},
		{
			name: "max countを指定する",
			opt: func(o *Opt[string]) {
				o.MaxCunt = 100
			},
			want: &writer[string]{
				maxCount:     100,
				emitDuration: 1 * time.Minute,
			},
		},
		{
			name: "emit durationを指定する",
			opt: func(o *Opt[string]) {
				o.EmitDuration = 100 * time.Second
			},
			want: &writer[string]{
				maxCount:     10,
				emitDuration: 100 * time.Second,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result Aggregator[string]
			if test.opt != nil {
				result = New(new(bytes.Buffer), test.opt)
			} else {
				result = New[string](new(bytes.Buffer))
			}

			w, ok := result.(*writer[string])
			if !ok {
				t.Fatal("要求する型ではない")
			}

			assert.NotNil(t, w.containers)
			assert.Equal(t, test.want.maxCount, w.maxCount)
			assert.Equal(t, test.want.emitDuration, w.emitDuration)
			assert.NotNil(t, w.aggregatedLogStrPipe)
			assert.NotNil(t, w.parse)
			assert.NotNil(t, w.keyGenerate)
		})
	}
}
