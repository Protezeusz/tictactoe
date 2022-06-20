package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"../pkg/files"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("failed read: %s", err)
			c.String(http.StatusBadRequest, "Invalid body.")
		}

		errIMg := files.SaveRequestImage(requestBody)
		if errIMg != nil {
			log.Fatalf("failed save: %s", err)
			c.String(http.StatusBadRequest, "Invalid body.")
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
