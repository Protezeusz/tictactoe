package files

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type IOFile struct {
	path string //  "rectangle.png"
}

func NewIOFile(p string) *IOFile {
	return &IOFile{
		path: p,
	}
}

func (file IOFile) SaveRequestImage(data []byte) error {

	coI := strings.Index(string(data), ",")
	rawImage := string(data)[coI+1:]

	unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))

	res := bytes.NewReader(unbased)

	pngI, errPng := png.Decode(res)
	if errPng != nil {
		log.Fatalf("failed encode request: %s", errPng)
		return errPng
	}

	f, err := os.Create(file.path)
	if err != nil {
		log.Fatalf("failed create file: %s", err)
		return err
	}
	png.Encode(f, pngI)
	return nil
}

func (file IOFile) ReadImage() ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(file.path)
	if err != nil {
		log.Fatalf("failed read file: %s", err)
		return nil, err
	}

	return fileBytes, nil
}
