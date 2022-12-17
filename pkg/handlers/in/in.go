package in

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/aggregator"
	"github.com/ieee0824/collei/pkg/handler"
	"github.com/ieee0824/collei/pkg/request/in"
)

func New(w io.Writer) *In {
	return &In{
		out: w,
		w:   make(map[string]io.Writer),
	}
}

type In struct {
	out io.Writer
	w   map[string]io.Writer
	handler.Handler
}

func (impl *In) Methods() []string {
	return []string{
		"POST",
	}
}

func (impl *In) Path() string {
	return "/in"
}

func sumSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (impl *In) Post(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	req, err := in.New(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "failed")
		return
	}

	_, ok := impl.w[req.Tag]
	if !ok {
		impl.w[req.Tag] = aggregator.New(impl.out, func(o *aggregator.Opt[map[string]any]) {
			o.MaxCunt = 3
			o.KeyGenerator = func(t *map[string]any) (string, error) {
				caller, ok := (*t)["caller"]
				if !ok {
					return sumSha1(fmt.Sprint(t)), nil
				}
				str, ok := caller.(string)
				if !ok {
					return sumSha1(fmt.Sprint(t)), nil
				}
				return sumSha1(str), nil
			}
		})
	}
	impl.w[req.Tag].Write(req.Body)
	ctx.JSON(http.StatusOK, "success")
}
