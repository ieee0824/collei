package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Path() string
	Methods() []string
	Post(*gin.Context)
}
