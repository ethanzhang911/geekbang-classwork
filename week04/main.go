package main

import (
	"context"
	"github.com/ethanzhang911/geekbang-classwork/v1/pkg/hstrix"
	"github.com/gin-gonic/gin"
)

var h hstrix.Hstrix

func pong(c *gin.Context) {
	if err := h.Add(); err != nil {
		c.String(400,"%v",err)
		return
	}else {
		c.String(200,"%s","pong")
	}

}

func main() {
	r := gin.Default()
	h = hstrix.NewHstrixByEthan(context.Background(), 4, 5, 1)
	h.Run()
	r.GET("/ping", pong)
	r.Run("127.0.0.1:8100")
}
