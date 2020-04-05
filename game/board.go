package game

import (
	"fmt"
	"math/rand"
)

const (
	blank =  0
	dugUp = -1
	mine  = -2
	flag  = -3
)

const (
	Ok = iota
	DigOnFlag = iota
	MineTriggered = iota
	OutOfBounds = iota
	AlreadyDugUp = iota
)

type Board struct {
	rows uint
	cols uint
	board [][]int
	spacesLeft uint
}

func NewBoard(m, rows, cols uint) Board {
	board := make([][]int, rows, rows)

	for i := uint(0); i < rows; i++ {
		board[i] = make([]int, cols, cols)
	}

	for currentMines := uint(0); currentMines < m; {
		r := rand.Uint32()

		i := uint(r & 0x000F) % rows
		j := uint(r & 0x00F0) % cols

		if board[i][j] == mine { continue }
		board[i][j] = mine
		currentMines++
	}

	return Board{
		rows: rows,
		cols: cols,
		board: board,
		spacesLeft: (rows * cols) - m,
	}
}

func (b *Board) Print() {
	fmt.Print("   ")
	for i := uint(0); i < b.cols; i++ {
		fmt.Printf(" %2d", i)
	}
	fmt.Println()

	for i, subarray := range b.board {
		fmt.Printf("%2d[", i)
		for _, val := range subarray {
			switch val {
			case blank, mine:
				fmt.Print("  █")
			case dugUp:
				fmt.Print("   ")
			case flag:
				fmt.Print("  ⚐")
			}
		}
		fmt.Println("  ]")
	}
}

func (b *Board) Dig(i, j uint) int {
	if i > b.rows - 1 || j > b.cols - 1 {
		return OutOfBounds
	}

	switch b.board[i][j] {
	case mine:
		return MineTriggered
	case flag:
		return DigOnFlag
	case blank:
		b.look(i, j)
		return Ok
	default:
		return AlreadyDugUp
	}
}

func (b *Board) look(i, j uint) {
	adjacentMines := 0

	type tuple struct { x, y uint }

	queue := make([]tuple, 0, 8)

	for wprobe := i - 1; wprobe <= i + 1 && wprobe >= 0 && wprobe < b.rows; wprobe++ {
		for hprobe := j - 1; hprobe <= j + 1 && hprobe >= 0 && hprobe < b.cols; hprobe++ {
			switch b.board[wprobe][hprobe] {
			case mine:
				adjacentMines++
			case blank:
				queue = append(queue, tuple{wprobe, hprobe})
			default:
				// no op; flag, dugUp or has adjacent mines already
			}
		}
	}

	if adjacentMines == 0 {
		b.board[i][j] = dugUp
	} else {
		b.board[i][j] = adjacentMines
	}

	for _, coordinate := range queue {
		b.look(coordinate.x, coordinate.y)
	}
}
