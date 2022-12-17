package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ieee0824/collei/pkg/apis"
)

func main() {
	r := gin.Default()
	f, e := os.Create("test.log")
	if e != nil {
		log.Fatalln(e)
	}
	defer f.Close()

	apis.New(f).RegistHandlers(r)

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
