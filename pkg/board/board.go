package board

import (
	"image"

	"github.com/Protezeusz/tictactoe/pkg/helper"
)

type Board struct {
	matrixboard [300][300]int64
	simpleBoard [3][3]string
}

func (b *Board) GetSimpleBoard() [3][3]string {
	return b.simpleBoard
}

func (b *Board) GetMatrixBoard() [300][300]int64 {
	return b.matrixboard
}

func (b *Board) UpdateBoard(img image.Image) error {

	var changed [3][3]int

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if pixel := helper.GetPixel(img.At(y, x)); b.matrixboard[y][x] != pixel {
				b.matrixboard[y][x] = pixel
				changed[y/100][x/100] += 1
			}
		}
	}

	return nil
}

func (b *Board) DrawXAt(x, y int) {
	// TODO
}

func (b *Board) DrawOAt(x, y int) {
	// TODO
}

func New() Board {
	return Board{
		matrixboard: [300][300]int64{},
		simpleBoard: [3][3]string{},
	}
}
