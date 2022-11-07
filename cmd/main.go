package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Protezeusz/tictactoe/pkg/board"
	"github.com/Protezeusz/tictactoe/pkg/helper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	board := board.New()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	r.POST("/", func(c *gin.Context) {

		requestBody, readErr := ioutil.ReadAll(c.Request.Body)
		if readErr != nil {
			log.Fatalf("failed: %s", readErr)
			c.String(http.StatusBadRequest, "Bad request.")
		}

		img, getImgErr := helper.GetImageFromHttp(requestBody)
		if getImgErr != nil {
			log.Fatalf("failed decode image request: %s", getImgErr)
			c.String(http.StatusInternalServerError, "Internal Server Error.")
		}

		updatErr := board.UpdateBoard(img)
		if updatErr != nil {
			log.Fatalf("failed: %s", updatErr)
			c.String(http.StatusInternalServerError, "Internal Server Error.")
		}

		// play game
		newImg, playErr := board.Play(img)
		if playErr != nil {
			log.Fatalf("failed: %s", updatErr)
			c.String(http.StatusInternalServerError, "Internal Server Error.")
		}

		// response
		byteImage, getErr := helper.GetHttpFromImage(newImg)
		if getErr != nil {
			log.Fatalf("failed: %s", getErr)
			c.String(http.StatusInternalServerError, "Internal Server Error.")
		}

		c.Header("Content-Disposition", "attachment; filename=file-name.txt")
		c.Data(http.StatusOK, "data:image/png;base64", byteImage)
		c.String(http.StatusOK, "ok")
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
