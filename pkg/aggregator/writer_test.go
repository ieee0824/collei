package aggregator

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name         string
		parser       ParserFunc[string]
		keyGenerator KeyGeneratorFunc[string]
		body         []byte
		isErr        bool
	}{
		{
			name: "parseできない",
			parser: func(b []byte) (*string, error) {
				return nil, errors.New("failed parse")
			},
			isErr: true,
			body:  []byte("foo\nbar\nbaz\n"),
		},
		{
			name: "keyを生成できない",
			parser: func(b []byte) (*string, error) {
				return lo.ToPtr(string(b)), nil
			},
			keyGenerator: func(a *string) (string, error) {
				return "", errors.New("failed gen key")
			},
			isErr: true,
			body:  []byte("foo\nbar\nbaz\n"),
		},
		{
			name: "containerに追加される",
			parser: func(b []byte) (*string, error) {
				return lo.ToPtr(string(b)), nil
			},
			keyGenerator: func(a *string) (string, error) {
				return "aaaa", nil
			},
			body: []byte("foo\nbar\nbaz\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pipe := make(chan string, 1024)
			w := &writer[string]{
				containers:           map[string]*container[string]{},
				parse:                test.parser,
				keyGenerate:          test.keyGenerator,
				maxCount:             10,
				emitDuration:         10 * time.Second,
				aggregatedLogStrPipe: pipe,
			}

			n, err := w.Write(test.body)
			if test.isErr {
				assert.NotNil(t, err)
				assert.Equal(t, 0, n)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, len(test.body), n)
		})
	}
}

func TestNewWriter(t *testing.T) {
	w := newWriter(new(bytes.Buffer), &Opt[string]{
		MaxCunt:      10,
		EmitDuration: 1 * time.Minute,
	})
	assert.NotNil(t, w.containers)
	assert.Equal(t, 10, w.maxCount)
	assert.Equal(t, 1*time.Minute, w.emitDuration)
	assert.NotNil(t, w.aggregatedLogStrPipe)
}
