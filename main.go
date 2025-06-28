package main

import (
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go h.run()

	router := gin.New()
	router.Use(cors.Default()) // allows requests from all origins

	// Serve the HTML chat interface from file
	router.GET("/", func(c *gin.Context) {
		c.File("chat.html")
	})

	// HEAD support for monitoring - returns headers without body
	router.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	// Also serve it at /chat for compatibility
	router.GET("/chat", func(c *gin.Context) {
		c.File("chat.html")
	})

	// WebSocket endpoint
	router.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		serveWs(c.Writer, c.Request, roomId)
	})

	router.Run(*addr)
}
