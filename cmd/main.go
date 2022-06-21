package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Protezeusz/tictactoe/pkg/files"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	file := files.NewIOFile("board.png")
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	r.POST("/", func(c *gin.Context) {
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("failed: %s", err)
			c.String(http.StatusBadRequest, "Invalid body.")
		}

		errIMg := file.SaveRequestImage(requestBody)
		if errIMg != nil {
			log.Fatalf("failed save: %s", err)
			c.String(http.StatusInternalServerError, "Internal server error.")
		}

		// play game

		// response
		c.File("rectangle.png")
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
