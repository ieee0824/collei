package option

import "github.com/ieee0824/collei/pkg/logs"

type Option struct {
	Protocol string
	Host     string
	Port     int
	Tag      logs.Tag
}
