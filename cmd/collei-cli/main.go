package main

import (
	"io"
	"os"
	"sync"

	"github.com/ieee0824/collei/pkg/collei/client"
	"github.com/rs/zerolog"
)

func main() {
	c := client.New()
	logger := zerolog.New(io.MultiWriter(c, os.Stdout)).With().Timestamp().Caller().Logger()

	var wg sync.WaitGroup
	semaphore := make(chan bool, 100)
	for i := 0; i < 65536; i++ {
		wg.Add(1)
		semaphore <- true

		go func(i int, wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			logger.Error().Msgf("%d: test", i)
		}(i, &wg)
	}
	wg.Wait()
}
