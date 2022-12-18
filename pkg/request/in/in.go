package in

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/logs"
)

type PostRequest struct {
	Tag  logs.Tag `json:"tag"`
	Body []byte   `json:"body"`
}

func New(ctx *gin.Context) (*PostRequest, error) {
	var ret PostRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
