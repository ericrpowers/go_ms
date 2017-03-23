package main

import (
	"fmt"
	"strconv"
)

type game struct {
	*Board
	turn int
}

// NewGame - constructor
func NewGame() {
	var answer string
	g := new(game)

	// Check if user or solver will play
	for !(answer == "u" || answer == "b") {
		answer = GetInput("Will the user or the bot play?(u/b) ")
	}

	if answer == "u" {
		g.User()
	} else {
		// TODO: Implement Solver
		// NewSolver()
	}
}

func (g *game) User() {
	// Initialize size and mines to use in loop
	var xSize, ySize, mines int
	var answer string
	var err error

	for i := 0; i < 3; i++ {
		tmp := -1
		fPass := true
		for tmp < 5 || (i <= 1 && tmp > 20) || (i > 1 && tmp > xSize*ySize/3) {
			if fPass {
				switch i {
				case 0:
					answer = GetInput("Size of X-axis (5 - 20): ")
				case 1:
					answer = GetInput("Size of Y-axis (5 - 20): ")
				case 2:
					/* Let's avoid too many mines and restrict it to
					   max ~1/3rd of entire board */
					answer = GetInput("Number of mines (5 - " +
						strconv.Itoa(xSize*ySize/3) + "): ")
				}
				fPass = false
			}

			// Look for an integer within the user input else ignore
			tmp, err = strconv.Atoi(answer)
			if err != nil {
				answer = ""
			}

			if tmp < 5 || (i <= 1 && tmp > 20) || (i > 1 && tmp > xSize*ySize/3) {
				dig := 20
				if i > 1 {
					dig = xSize * ySize / 3
				}
				answer = GetInput("Choose a number between 5 and " +
					strconv.Itoa(dig) + ": ")
			}
		}

		switch i {
		case 0:
			xSize = tmp
		case 1:
			ySize = tmp
		case 2:
			mines = tmp
		}
	}

	answer = "y"
	for answer == "y" {
		g.Board = NewBoard(xSize, ySize, mines)
		g.turn = 0
		g.Play()

		answer = ""
		for !(answer == "y" || answer == "n") {
			answer = GetInput("Want to play again?(y/n) ")
		}
	}
}

func (g *game) Play() {
	for {
		g.turn++
		fmt.Printf("\nTurn %d\n", g.turn)
		g.Board.ShowBoard()
		if g.Board.IsFinalMove(g.Board.SetPosition()) {
			break
		}
	}

	if g.Board.Win() {
		fmt.Printf("Congrats, you found all the mines in %d turns.\n", g.turn)
	} else {
		fmt.Println("You hit a mine! Try again.")
	}
	g.Board.ShowMines()
}
