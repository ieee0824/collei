package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ieee0824/collei/pkg/aggregator"
)

type data struct {
	Msg  string    `json:"msg,omitempty"`
	Date time.Time `json:"date,omitempty"`
}

func main() {
	a := aggregator.New[data](os.Stdout, func(o *aggregator.Opt[data]) {
		o.KeyGenerator = func(t *data) (string, error) {
			return t.Date.String(), nil
		}
		o.MaxCunt = 5
		o.EmitDuration = 1 * time.Second
	})

	now := time.Now()
	for i := 0; i < 4; i++ {
		err := json.NewEncoder(a).Encode(&data{
			Msg:  fmt.Sprintf("msg: %d", i),
			Date: now,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	time.Sleep(10 * time.Second)
}
