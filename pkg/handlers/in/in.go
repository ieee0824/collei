package in

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/aggregator"
	"github.com/ieee0824/collei/pkg/handler"
	"github.com/ieee0824/collei/pkg/request/in"
	"github.com/samber/lo"
)

type Opt struct {
	EmitCount    int
	EmitDuration time.Duration
}

type OptFunc func(opt *Opt)

func New(w io.Writer, of ...OptFunc) *In {
	opt := &Opt{
		EmitCount:    100,
		EmitDuration: 60 * time.Second,
	}
	lo.ForEach(of, func(f OptFunc, _ int) {
		f(opt)
	})
	return &In{
		emitCount:    opt.EmitCount,
		emitDuration: opt.EmitDuration,
		out:          w,
		w:            make(map[string]io.Writer),
	}
}

type In struct {
	emitCount    int
	emitDuration time.Duration
	out          io.Writer
	w            map[string]io.Writer
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
			o.MaxCnt = impl.emitCount
			o.EmitDuration = impl.emitDuration
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
