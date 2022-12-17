package aggregator

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name         string
		parser       ParserFunc[any]
		keyGenerator KeyGeneratorFunc[any]
		body         []byte
		isErr        bool
	}{
		{
			name: "parseできない",
			parser: func(b []byte) (any, error) {
				return nil, errors.New("failed parse")
			},
			isErr: true,
			body:  []byte("foo\nbar\nbaz\n"),
		},
		{
			name: "keyを生成できない",
			parser: func(b []byte) (any, error) {
				return strings.Split(string(b), "\n"), nil
			},
			keyGenerator: func(a any) (string, error) {
				return "", errors.New("failed gen key")
			},
			isErr: true,
			body:  []byte("foo\nbar\nbaz\n"),
		},
		{
			name: "containerに追加される",
			parser: func(b []byte) (any, error) {
				return strings.Split(string(b), "\n"), nil
			},
			keyGenerator: func(a any) (string, error) {
				return "aaaa", nil
			},
			body: []byte("foo\nbar\nbaz\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pipe := make(chan string, 1024)
			w := &writer[any]{
				containers:           map[string]*container[any]{},
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
