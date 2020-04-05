package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AnthonyHewins/minesweeper/game"
)

func scanCmd() {
	fmt.Println("Use (d uint uint) to dig, (f uint uint) to flag, exit strings to quit")
	var cmd string
	fmt.Scan(&cmd)

	switch cmd {
	case "e", "q", "exit", "quit":
		fmt.Println("bye")
		os.Exit(0)
	default:
	}
}

func main() {
	mines  := flag.Uint("mines",  10, "Number of mines, 1 to 99")
	cols  := flag.Uint("cols",  30, "Board cols, 1 to 99")
	rows := flag.Uint("rows", 30, "Board rows, 1 to 99")

	flag.Parse()

	check := func(ptr *uint, name string) {
		if *ptr >= 100 || *ptr < 1 {
			fmt.Println("%s was not between 1 and 99: got %s", name, *ptr)
			os.Exit(1)
		}
	}

	check(mines, "mines")
	check(cols, "cols")
	check(rows, "rows")

	if *mines >= (*cols) * (*rows) {
		fmt.Println("more/equal mines to spaces!")
		os.Exit(1)
	}

	board := game.NewBoard(*mines, *rows, *cols)

	board.Print()
}
