package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ieee0824/collei/pkg/collei/option"
	"github.com/ieee0824/collei/pkg/request/in"
	"github.com/samber/lo"
)

type optFunc func(opt *option.Option)

func New(optFuncs ...optFunc) Collie {
	defaultOption := &option.Option{
		Protocol: "http",
		Host:     "localhost",
		Port:     8080,
	}

	lo.ForEach(optFuncs, func(f optFunc, _ int) {
		f(defaultOption)
	})

	return &client{
		endpoint: genEndpoint(defaultOption),
		tag:      defaultOption.Tag,
	}
}

func genEndpoint(opt *option.Option) string {
	return fmt.Sprintf("%s://%s:%d", opt.Protocol, opt.Host, opt.Port)
}

type Collie interface {
	io.Writer
}

type client struct {
	endpoint string
	tag      string
}

func (c *client) Write(b []byte) (int, error) {
	req := &in.PostRequest{
		Tag:  c.tag,
		Body: b,
	}
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return 0, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/in", c.endpoint),
		body,
	)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	default:
		return 0, errors.New(resp.Status)
	}

	return len(b), nil
}
