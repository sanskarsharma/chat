package main

import (
	"flag"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

)
var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go h.run()

	router := gin.New()
	router.Use(cors.Default()) // allows requests from all origins 

	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		serveWs(c.Writer, c.Request, roomId)
	})

	router.Run(*addr)
}
