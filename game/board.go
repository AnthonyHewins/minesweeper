package game

import (
	"fmt"
	"time"
	"math/rand"
)

const (
	blank       =  0
	dugUp       = -1
	mine        = -2
	flag        = -3
	flaggedMine = -4
)

const (
	Ok = iota
	DigOnFlag = iota
	MineTriggered = iota
	OutOfBounds = iota
	AlreadyDugUp = iota
	Victory = iota
)

type tuple struct {
	x, y int
}

type Board struct {
	rows int
	cols int
	board [][]int8
	spacesLeft int
}

func NewBoard(m, rows, cols int) Board {
	board := make([][]int8, rows, rows)

	for i := 0; i < rows; i++ {
		board[i] = make([]int8, cols, cols)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r  := rand.New(s1)

	for currentMines := 0; currentMines < m; {
		i := int8(r.Int31n(int32(rows)))
		j := int8(r.Int31n(int32(cols)))

		if i < 0 { i *= -1 }
		if j < 0 { j *= -1 }

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

func (b *Board) Winrar() {
	b.printCols()
	for i, subarray := range b.board {
		fmt.Printf("%2d[", i)
		for _, val := range subarray {
			switch val {
			case mine:
				fmt.Print("  ☠")
			default:
				fmt.Print("   ")
			}
		}
		fmt.Printf("  ]%2d\n", i)
	}
	b.printCols()
}

func (b *Board) Print() {
	b.printCols()
	for i, subarray := range b.board {
		fmt.Printf("%2d[", i)
		for _, val := range subarray {
			switch val {
			case blank, mine:
				fmt.Print(" ██")
			case dugUp:
				fmt.Print("   ")
			case flag, flaggedMine:
				fmt.Print("  ⚐")
			default:
				fmt.Printf("  %d", val)
			}
		}
		fmt.Printf("  ]%2d\n", i)
	}
	b.printCols()
}

func (b *Board) printCols() {
	fmt.Print("   ")
	for i := 0; i < b.cols; i++ {
		fmt.Printf(" %2d", i)
	}
	fmt.Println()
}

func (b *Board) Spaces() int {
	return b.spacesLeft
}

func (b *Board) Dig(i, j int) int {
	if i > b.rows - 1 || j > b.cols - 1 {
		return OutOfBounds
	}

	switch b.board[i][j] {
	case mine:
		return MineTriggered
	case flag, flaggedMine:
		return DigOnFlag
	case blank:
		b.look(i, j)
		if b.spacesLeft == 0 { return Victory }
		return Ok
	default:
		return AlreadyDugUp
	}
}

func (b *Board) FlagToggle(i, j int) int {
	if i > b.rows - 1 || j > b.cols - 1 {
		return OutOfBounds
	}

	switch b.board[i][j] {
	case flag:
		b.board[i][j] = blank
	case blank:
		b.board[i][j] = flag
	case flaggedMine:
		b.board[i][j] = mine
	case mine:
		b.board[i][j] = flaggedMine
	default:
		// noop for revealed spaces; numbered and dugup
	}

	return Ok
}

func (b *Board) look(i, j int) {
	switch b.board[i][j] {
	case dugUp, mine, flaggedMine:
		return
	case blank, flag:
		// noop; if blank, we want to continue, find how many mines are around it
		// if flag, then the user was wrong about the placement, so we dig it up anyways
	default:
		return // already numbered with
	}

	b.spacesLeft--
	if b.spacesLeft == 0 { return }

	adjacentMines := int8(0)

	queue := make([]tuple, 0, 8)

	rowStart, rowEnd, colStart, colEnd := b.boundaryCheck(i - 1, i + 1, j - 1, j + 1)

	for wprobe := rowStart; wprobe <= rowEnd; wprobe++ {
		for hprobe := colStart; hprobe <= colEnd; hprobe++ {
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
		for _, coordinate := range queue {
			b.look(coordinate.x, coordinate.y)
		}
	} else {
		b.board[i][j] = adjacentMines
	}
}

func (b *Board) boundaryCheck(rowStart, rowEnd, colStart, colEnd int) (int, int, int, int) {
	if rowStart < 0 {
		rowStart = 0
	}

	if colStart < 0 {
		colStart = 0
	}

	if rowEnd >= b.rows {
		rowEnd = b.rows - 1
	}

	if colEnd >= b.cols {
		colEnd = b.cols - 1
	}

	return rowStart, rowEnd, colStart, colEnd
}
