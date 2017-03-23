package main

import (
	"bytes"
	"fmt"
	"strconv"
)

// This mapping will be how the field is displayed
const coverSymbol, emptySymbol, mineSymbol = "X", ".", "*"

// Board - structure
type Board struct {
	colLen, rowLen, mCount, row, column int
	boardgame                           []int
	*Minefield
}

// NewBoard - constructor
func NewBoard(xSize, ySize, mines int) *Board {
	b := new(Board)
	b.colLen = xSize
	b.rowLen = ySize
	b.mCount = mines
	b.row, b.column = -1, -1
	// Only setup the board to avoid losing on the first turn
	b.startBoard()

	return b
}

func (b *Board) startBoard() {
	b.boardgame = make([]int, b.colLen*b.rowLen)
	for i := 0; i < len(b.boardgame); i++ {
		b.boardgame[i] = 9
	}
}

// GetBoardValues - Retrieves entire boardgame array
func (b *Board) GetBoardValues() []int {
	return b.boardgame
}

// IsCellCovered - Checks whether or not a cell is in a covered state
func (b *Board) IsCellCovered(c, r int) bool {
	return b.boardgame[(b.colLen*r)+c] == 9
}

// ShowBoard - Shows the current state of the board in ASCII form
func (b *Board) ShowBoard() {
	// TODO: Cleanup layout to make things prettier
	fmt.Printf("%11s", "Rows\n")
	for x := b.rowLen; x > 0; x-- {
		fmt.Printf("%9d ", x)

		for y := 1; y <= b.rowLen; y++ {
			var rLine bytes.Buffer
			rLine.WriteString("    ")
			val := b.boardgame[(b.colLen*(x-1))+(y-1)]
			switch val {
			case -1:
				rLine.WriteString(mineSymbol)
			case 0:
				rLine.WriteString(emptySymbol)
			case 9:
				rLine.WriteString(coverSymbol)
			default:
				rLine.WriteString(strconv.Itoa(val))
			}
			fmt.Print(rLine.String())
		}
		fmt.Println()
	}
	var cLine bytes.Buffer
	cLine.WriteString("\n              1")
	for i := 2; i <= b.colLen; i++ {
		if i < 10 {
			cLine.WriteString("    " + strconv.Itoa(i))
		} else {
			cLine.WriteString("   " + strconv.Itoa(i))
		}
	}
	fmt.Println(cLine.String())
	fmt.Printf(("%" + strconv.Itoa(cLine.Len()/2+9) + "s"), "Columns\n")
}

// SetPosition - Sets the next position is be discovered
func (b *Board) SetPosition() bool {
	for {
		var itr, dig int
		var answer string
		var err error
		b.row = -1
		b.column = -1
		// Do not leave this region until two valid inputs are given
		for (itr == 0 && (b.row < 0 || b.row > b.rowLen-1)) ||
			(itr != 0 && (b.column < 0 || b.column > b.colLen-1)) {
			if itr == 0 {
				answer = GetInput("Row: ")
			} else {
				answer = GetInput("Column: ")
			}

			if itr == 0 {
				b.row, err = strconv.Atoi(answer)
				b.row--
			} else {
				b.column, err = strconv.Atoi(answer)
				b.column--
			}

			dig = b.colLen
			if itr == 0 {
				dig = b.rowLen
			}
			if err != nil {
				fmt.Printf("Choose a number between 1 and %d\n", dig)
				continue
			}
			if (itr == 0 && (b.row < 0 || b.row > b.rowLen-1)) ||
				(itr != 0 && (b.column < 0 || b.column > b.colLen-1)) {
				fmt.Printf("Choose a number between 1 and %d\n", dig)
			} else {
				itr++
			}
		}

		if !b.IsCellCovered(b.column, b.row) {
			fmt.Println("Field already shown")
		} else {
			break
		}
	}

	// Return true if a mine is hit
	return b.GetPositionVal(b.column, b.row) == b.Minefield.GetMineVal()
}

// GetPositionVal - Retrieves the positions value (empty, hint, or mine)
func (b *Board) GetPositionVal(y, x int) int {
	// Setup the minefield if it does not exist
	b.column = y
	b.row = x
	if b.Minefield == nil {
		b.Minefield = NewMineField(b.colLen, b.rowLen, b.mCount, x, y)
	}

	return b.Minefield.GetCellVal((b.colLen * x) + y)
}

// IsFinalMove - Finds out if the game is over
func (b *Board) IsFinalMove(isMine bool) bool {
	// If user did not hit a mine, uncover the empty region
	if !isMine {
		b.openNeighbors()
		isMine = b.Win()
	}

	return isMine
}

/*  When revealing the playing field, need to take into account the boundaries
    of the playing field and if the position is an empty cell (aka 0). If it is
    an empty cell, recursively continue to expand range until the entire empty
    region is exposed. */
func (b *Board) openNeighbors() {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if b.row+i < 0 || b.row+i >= b.rowLen || b.column+j < 0 ||
				b.column+j >= b.colLen {
				continue
			} else if b.Minefield.GetCellVal((b.colLen*(b.row+i))+(b.column+j)) ==
				b.Minefield.GetMineVal() {
				continue
			} else if !b.IsCellCovered(b.column+j, b.row+i) {
				continue
			}

			val := b.Minefield.GetCellVal((b.colLen * (b.row + i)) + (b.column + j))
			b.boardgame[b.colLen*(b.row+i)+(b.column+j)] = val
			if val == 0 && !(b.row+i == b.row && b.column+j == b.column) {
				b.row += i
				b.column += j
				b.openNeighbors()
				b.column -= j
				b.row -= i
			}
		}
	}
}

// Win - Determines if the board was successfully cleared
func (b *Board) Win() bool {
	var count int
	for x := 0; x < b.rowLen; x++ {
		for y := 0; y < b.colLen; y++ {
			if b.IsCellCovered(y, x) {
				count++
			}
		}
	}

	return count == b.mCount
}

// ShowMines - Shows the mine locations on the board
func (b *Board) ShowMines() {
	for x := 0; x < b.rowLen; x++ {
		for y := 0; y < b.colLen; y++ {
			if b.Minefield.GetCellVal((b.colLen*x)+y) ==
				b.Minefield.GetMineVal() {
				b.boardgame[(b.colLen*x)+y] = b.Minefield.GetMineVal()
			}
		}
	}

	b.ShowBoard()
}
