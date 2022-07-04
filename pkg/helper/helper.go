package helper

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

func GetPixel(c color.Color) int64 {
	white := color.NRGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}

	if c == white {
		return 0
	}
	return 1
}

func GetImageFromHttp(requestBody []byte) (image.Image, error) {

	coI := strings.Index(string(requestBody), ",")
	rawImage := string(requestBody)[coI+1:]

	unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))

	res := bytes.NewReader(unbased)

	imgByte, decodErr := png.Decode(res)
	if decodErr != nil {
		return nil, decodErr
	}

	return imgByte, nil

}

func GetHttpFromImage(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	encodErr := png.Encode(buf, img)
	if encodErr != nil {
		return nil, encodErr
	}

	return buf.Bytes(), nil
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
