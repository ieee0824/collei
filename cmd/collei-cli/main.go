package main

import (
	"io"
	"os"

	"github.com/ieee0824/collei/pkg/collei/client"
	"github.com/rs/zerolog"
)

func main() {
	c := client.New()
	logger := zerolog.New(io.MultiWriter(c, os.Stdout)).With().Timestamp().Caller().Logger()

	logger.Error().Msg("test")
}
