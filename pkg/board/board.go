package board

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Protezeusz/tictactoe/pkg/helper"
	"github.com/montanaflynn/stats"
)

type Board struct {
	matrixboard [300][300]float64
	simpleBoard [3][3]string
	playAs      string
	winer       [3]stats.Coordinate
}

type move struct {
	x int
	y int
}

func (b *Board) GetSimpleBoard() [3][3]string {
	return b.simpleBoard
}

func (b *Board) GetMatrixBoard() [300][300]float64 {
	return b.matrixboard
}

func (b *Board) GetPlayAs() string {
	return b.playAs
}

func (b *Board) UpdateBoard(img image.Image) error {

	var points []stats.Coordinate
	var sumX float64 = 0
	var sumY float64 = 0

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if pixel := helper.GetPixel(img.At(y, x)); b.matrixboard[y][x] != pixel {
				b.matrixboard[y][x] = pixel
				points = append(points, stats.Coordinate{X: float64(x), Y: float64(y)})
				sumX += float64(x)
				sumY += float64(y)
			}
		}
	}

	centerX := sumX / float64(len(points))
	centerY := sumY / float64(len(points))
	fmt.Printf("X: %f\nY: %f\n", centerX, centerY)

	result, checkErr := helper.CheckIfXOrY(centerX, centerY, points)
	if checkErr != nil {
		return checkErr
	}

	b.simpleBoard[int(centerY)/100][int(centerX)/100] = result
	switch result {
	case "O":
		b.playAs = "X"
	case "X":
		b.playAs = "O"
	default:
		return fmt.Errorf("unknown player")
	}
	return nil
}

func (b *Board) print() {
	for x := 0; x < 3; x++ {
		fmt.Print("[")
		for y := 0; y < 3; y++ {
			if b.simpleBoard[y][x] != "" {
				fmt.Printf(" %s", b.simpleBoard[y][x])
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print(" ]\n")
	}
}
func (b *Board) Play(img image.Image) (image.Image, error) {

	score := b.evaluate(b.simpleBoard)
	if score != 0 {
		return b.drawEnd(img, score)
	}

	nextMove := b.findBestMove()
	b.print()

	if nextMove.x != -1 {
		switch b.playAs {
		case "X":
			return b.drawXAt(nextMove.y, nextMove.x, img)
		case "O":
			return b.drawOAt(nextMove.y, nextMove.x, img)
		default:
			return img, fmt.Errorf("unknown player")
		}
	}

	score = b.evaluate(b.simpleBoard)
	return b.drawEnd(img, score)
}

func (b *Board) findBestMove() move {
	bestVal := math.Inf(-1)
	bestMove := move{x: -1, y: -1}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if b.simpleBoard[y][x] == "" {
				b.simpleBoard[y][x] = b.playAs
				movetVal := b.minimax(b.simpleBoard, 0, false)
				b.simpleBoard[y][x] = ""
				if movetVal > bestVal {
					bestMove.x = x
					bestMove.y = y
					bestVal = movetVal
				}
			}
		}
	}
	if bestVal != math.Inf(-1) {
		b.simpleBoard[bestMove.y][bestMove.x] = b.playAs
	}

	return bestMove
}

func isMovesLeft(board [3][3]string) bool {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if board[y][x] == "" {
				return true
			}
		}
	}
	return false
}

func (b *Board) evaluate(board [3][3]string) float64 {
	for row := 0; row < 3; row++ {
		if board[row][0] == board[row][1] && board[row][1] == board[row][2] {
			if board[row][0] != "" {
				for i := 0; i < 3; i++ {
					b.winer[i] = stats.Coordinate{
						X: float64(50 + i*100),
						Y: float64(50 + row*100),
					}
				}
				if board[row][0] == b.playAs {
					return 10
				} else {
					return -10
				}
			}
		}
	}
	for col := 0; col < 3; col++ {
		if board[0][col] == board[1][col] && board[1][col] == board[2][col] {
			if board[0][col] != "" {
				for i := 0; i < 3; i++ {
					b.winer[i] = stats.Coordinate{
						X: float64(50 + col*100),
						Y: float64(50 + i*100),
					}
				}
				if board[0][col] == b.playAs {
					return 10
				} else {
					return -10
				}
			}
		}
	}
	if board[1][1] != "" {
		if board[0][0] == board[1][1] && board[1][1] == board[2][2] {
			for i := 0; i < 3; i++ {
				b.winer[i] = stats.Coordinate{
					X: float64(50 + i*100),
					Y: float64(50 + i*100),
				}
			}
			if board[1][1] == b.playAs {
				return 10
			} else {
				return -10
			}
		}
		if board[0][2] == board[1][1] && board[1][1] == board[2][0] {
			for i := 0; i < 3; i++ {
				b.winer[i] = stats.Coordinate{
					X: float64(50 + i*100),
					Y: float64(50 + (2-i)*100),
				}
			}
			if board[1][1] == b.playAs {
				return 10
			} else {
				return -10
			}
		}
	}
	return 0
}

func (b *Board) minimax(board [3][3]string, depth int, isMaximizingPlayer bool) float64 {
	score := b.evaluate(board)

	if score == 10 {
		return score - float64(depth)
	}
	if score == -10 {
		return score + float64(depth)
	}
	if isMovesLeft(board) {
		return 0
	}

	if isMaximizingPlayer {
		// X, Y
		best := math.Inf(-1)
		for y := 0; y < 3; y++ {
			for x := 0; x < 3; x++ {
				if board[y][x] == "" {
					board[y][x] = b.playAs
					value := b.minimax(board, depth+1, false)
					best = math.Max(best, value)
					board[y][x] = ""
				}
			}
		}
		return best
	} else {
		best := math.Inf(1)
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				if board[y][x] == "" {
					board[y][x] = "*"
					value := b.minimax(board, depth+1, true)
					best = math.Min(best, value)
					board[y][x] = ""
				}
			}
		}
		return best
	}
}

func (b *Board) drawEnd(img image.Image, score float64) (image.Image, error) {

	if score != 0 {
		A := b.winer[0].Y - b.winer[2].Y
		B := b.winer[2].X - b.winer[0].X
		C := A*b.winer[2].X + B*b.winer[2].Y
		for i := 0; i < 3; i++ {
			x := b.winer[i].X - 50
			y := b.winer[i].Y - 50
			for i := y; i < 100+y; i++ {
				for j := x; j < 100+x; j++ {
					if C-1000 < A*j+B*i && C+1000 > A*j+B*i {
						b.matrixboard[int(i)][int(j)] = 1
						if cimg, ok := img.(helper.Changeable); ok {
							if score < 0 {
								cimg.Set(int(i), int(j), color.RGBA{
									R: 0,
									G: 255,
									B: 0,
									A: 255,
								})
							} else {
								cimg.Set(int(i), int(j), color.RGBA{
									R: 255,
									G: 0,
									B: 0,
									A: 255,
								})
							}
						} else {
							return img, errors.New("image not changeable")
						}

					}
				}
			}
		}
	}

	return img, nil
}

func (b *Board) drawXAt(y, x int, img image.Image) (image.Image, error) {
	for i := y * 100; i < 100+y*100; i++ {
		for j := x * 100; j < 100+x*100; j++ {
			q := i - y*100
			w := j - x*100
			if w-5 < q && q < w+5 || 95-w < q && q < 105-w {
				b.matrixboard[i][j] = 1
				if cimg, ok := img.(helper.Changeable); ok {
					cimg.Set(i, j, color.Black)
				} else {
					return img, errors.New("image not changeable")
				}

			}
		}
	}

	if score := b.evaluate(b.simpleBoard); score != 0 {
		return b.drawEnd(img, score)
	}
	return img, nil
}

func (b *Board) drawOAt(y, x int, img image.Image) (image.Image, error) {
	for i := y * 100; i < 100+y*100; i++ {
		for j := x * 100; j < 100+x*100; j++ {
			var q float64 = float64(i - y*100)
			var w float64 = float64(j - x*100)
			f := math.Pow(q-50, 2) + math.Pow(w-50, 2)
			if 2000 < f && f < 2500 {
				b.matrixboard[i][j] = 1
				if cimg, ok := img.(helper.Changeable); ok {
					cimg.Set(i, j, color.Black)
				} else {
					return img, errors.New("image not changeable")
				}
			}
		}
	}
	if score := b.evaluate(b.simpleBoard); score != 0 {
		return b.drawEnd(img, score)
	}
	return img, nil
}

func New() Board {
	return Board{
		matrixboard: [300][300]float64{},
		simpleBoard: [3][3]string{},
	}
}
