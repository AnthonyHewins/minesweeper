package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AnthonyHewins/minesweeper/game"
)

type command struct {
	command string
	x       int
	y       int
}

func scanCmd(spacesLeft int) command {
	for {
		fmt.Printf("%d space(s) left. Use (d uint uint) to dig, (f uint uint) to flag, exit/quit/e/q to quit\n", spacesLeft)
		var cmd string

		_, err := fmt.Scan(&cmd)
		if err != nil {
			fmt.Println("error reading first arg")
			continue
		}

		switch cmd {
		case "e", "exit", "quit", "q":
			fmt.Println("bye")
			os.Exit(0)
		}

		var x, y int
		if _, err := fmt.Scanf("%d", &x); err != nil {
			fmt.Println("error reading x")
		}

		if _, err := fmt.Scanf("%d", &y); err != nil {
			fmt.Println("error reading y")
		}
	
		return command{command: cmd, x: x, y: y}
	}
}

func flagResult(result int) {
	switch result {
	case game.OutOfBounds:
		fmt.Println("coordinates given are out of bounds")
	}
}

func digResult(result int) bool {
	switch result {
	case game.DigOnFlag:
		fmt.Println("Can't dig on a flag")
	case game.MineTriggered:
		fmt.Println("You died, you suck")
		os.Exit(0)
	case game.OutOfBounds:
		fmt.Println("coordinates given are out of bounds")
	case game.AlreadyDugUp, game.Ok:
		// noop
	case game.Victory:
		return true
	default:
		fmt.Println("unknown result from API")
		os.Exit(1)
	}

	return false
}

func main() {
	mines := flag.Int("mines", 10, "Number of mines, 1 to 99")
	rows  := flag.Int("rows",  30, "Board rows, 1 to 99")
	cols  := flag.Int("cols",  30, "Board cols, 1 to 99")

	flag.Parse()

	check := func(ptr *int, name string) {
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

	fmt.Printf("Starting %d x %d game with %d mines\n", *rows, *cols, *mines)

	for {
		board.Print()
		cmd := scanCmd(board.Spaces())

		switch cmd.command {
		case "e", "q", "exit", "quit":
			fmt.Println("bye")
			os.Exit(0)
		case "d", "dig":
			if digResult(board.Dig(cmd.y, cmd.x)) {
				board.Winrar()
				fmt.Println("you won")
				os.Exit(0)
			}
		case "f", "flag", "rf", "remove-flag":
			flagResult(board.FlagToggle(cmd.x, cmd.y))
		default:
			fmt.Println("can't understand command: %s", cmd.command)
		}
	}
}
