package helper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
)

type Changeable interface {
	Set(x, y int, c color.Color)
}

func GetPixel(c color.Color) float64 {
	white := color.NRGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}

	if c == white {
		return 0
	}
	return 1.0
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

func GetDet(matrix [][]float64) (float64, error) {
	n := len(matrix)
	if n != len(matrix[0]) {
		return 0.0, errors.New("matrix is not square")
	}
	if n < 1 {
		return 0.0, errors.New("n < 1")
	} else {
		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		sd := 1.0
		for k := 0; k < n; k++ {
			max := math.Abs(matrix[k][k])
			ih := k
			jh := k
			for i := k; i < n; i++ {
				for j := k; j < n; j++ {
					if k == 0 {
						matrix[i][j] = matrix[i][j] * random.Float64() * 100
					}
					m := math.Abs(matrix[i][j])
					if m > max {
						max = m
						ih = i
						jh = j
					}
				}
			}
			if max == 0 {
				return max, nil
			} else {
				if ih != k {
					sd = -sd
					for j := k; j < n; j++ {
						m := matrix[k][j]
						matrix[k][j] = matrix[ih][j]
						matrix[ih][j] = m
					}
				}
				if jh != k {
					sd = -sd
					for i := k; i < n; i++ {
						m := matrix[i][k]
						matrix[i][k] = matrix[i][jh]
						matrix[i][jh] = m
					}
				}
				for i := k + 1; i < n; i++ {
					m := matrix[i][k] / matrix[k][k]
					for j := k; j < n; j++ {
						matrix[i][j] -= m * matrix[k][j]
					}
				}
			}
		}
		if matrix[n-1][n-1] == 0 {
			return matrix[n-1][n-1], nil
		} else {
			m := matrix[0][0]
			for i := 1; i < n; i++ {
				m *= matrix[i][i]
			}
			m *= sd
			return m, nil
		}
	}
}

func CheckIfXOrY(X, Y float64, points []stats.Coordinate) (string, error) {

	radiuses := make([]float64, len(points))
	for i := 0; i < len(points); i++ {
		radiuses[i] = math.Sqrt(math.Pow(X-points[i].X, 2) + math.Pow(Y-points[i].Y, 2))
	}

	deviation, statErr := stats.StandardDeviation(radiuses)
	if statErr != nil {
		return "", statErr
	}

	meanRadius, meanErr := stats.Mean(radiuses)
	if meanErr != nil {
		return "", meanErr
	}

	fmt.Printf("deviation: %f\n", deviation)
	fmt.Printf("meanRadius: %f\n", meanRadius)
	if deviation < 0.4*meanRadius {
		fmt.Println("koło")
		return "O", nil
	} else {
		fmt.Println("krzyżyk")
		return "X", nil
	}
}
