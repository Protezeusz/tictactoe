package files

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"log"
	"os"
	"strings"
)

func SaveRequestImage(data []byte) error {

	coI := strings.Index(string(data), ",")
	rawImage := string(data)[coI+1:]

	unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))

	res := bytes.NewReader(unbased)

	pngI, errPng := png.Decode(res)
	if errPng != nil {
		return errPng
	}

	rectangle := "rectangle.png"

	file, err := os.Create(rectangle)
	if err != nil {
		log.Fatalf("failed create file: %s", err)
	}
	png.Encode(file, pngI)
	return nil
}
