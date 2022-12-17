package apis

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/handler"
	"github.com/ieee0824/collei/pkg/handlers/in"
	"github.com/samber/lo"
)

func New(w io.Writer) *APIs {
	return &APIs{
		handlers: []handler.Handler{
			in.New(w),
		},
	}
}

type APIs struct {
	handlers []handler.Handler
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
