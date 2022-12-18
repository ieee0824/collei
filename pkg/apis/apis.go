package apis

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/handler"
	"github.com/ieee0824/collei/pkg/handlers/in"
	"github.com/ieee0824/collei/pkg/logs"
	"github.com/samber/lo"
)

func New(w io.Writer) *APIs {
	ls := logs.Logs{}
	return &APIs{
		handlers: []handler.Handler{
			in.New(w, func(opt *in.Opt) {
				opt.Logs = ls
			}),
		},
		logs: ls,
	}
}

type APIs struct {
	handlers []handler.Handler
	logs     logs.Logs
}

func (impl *APIs) RegistHandlers(engine *gin.Engine) {
	lo.ForEach(impl.handlers, func(h handler.Handler, _ int) {
		lo.ForEach(h.Methods(), func(method string, _ int) {
			switch method {
			case "POST":
				engine.POST(h.Path(), h.Post)
			default:
				log.Printf("unsupported http method: %s", method)
			}
		})
	})
}
