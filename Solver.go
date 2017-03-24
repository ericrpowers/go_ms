package main

import (
	"fmt"
	//"log"
	"time"
)

type solver struct {
	xSize, ySize, mines, numOfGames, wins, safeMoves, numOfTurns int
	savedBoard                                                   []int
	cellArray                                                    []Cell
	*Board
}

// NewSolver - This will house the strategy/logic needed to solve Minesweeper
func NewSolver() {
	/* The preset size of the board will be 10x10
	   The preset number of mines will be 10 */
	s := new(solver)
	s.xSize, s.ySize, s.mines = 10, 10, 10
	s.numOfGames, s.wins, s.safeMoves = 0, 0, 0
	s.numOfTurns = 1

	// 10,000 is a good statistical average
	iterations := 100
	startTime := time.Now()

	for iterations != 0 {
		s.numOfGames++
		s.Board = NewBoard(s.xSize, s.ySize, s.mines)
		s.savedBoard = nil
		coordArray := s.SelectNextCell()
		for !s.IsFinalMove(s.GetPositionVal(coordArray[0], coordArray[1]) == -1) {
			s.numOfTurns++
			// s.ShowBoard()
			coordArray = s.SelectNextCell()
		}

		if s.Win() {
			s.wins++
		} //else {
		//	s.ShowBoard()
		//	log.Fatal("Debug: Lost game")
		//}
		iterations--
	}
	elapsedTime := time.Since(startTime)

	// Now to print out the statistics
	fmt.Printf("%20s %10d\n", "Total games", s.numOfGames)
	fmt.Printf("%20s %10d\n", "Total wins", s.wins)
	fmt.Printf("%20s %10d\n", "Total time taken (s)", elapsedTime)
	fmt.Printf("%20s %10d\n", "Win Percentage (%)", s.wins/(s.numOfGames*1.00)*100)
	fmt.Printf("%20s %10d\n", "Average # of turns", s.numOfTurns/(s.numOfGames*1.00))
}

// SelectNextCell - Blind first move, then safe moves, and finally risky moves
func (s *solver) SelectNextCell() []int {
	var arry []int
	if s.savedBoard == nil {
		s.safeMoves = 0
		s.savedBoard = []int{s.xSize / 2, s.ySize / 2}

		// First round just click in the center, since we know we cannot lose
		return s.savedBoard
	}

	if s.safeMoves == 0 {
		s.savedBoard = s.GetBoardValues()
		s.cellArray = make([]Cell, s.xSize*s.ySize)
		for i := 0; i < s.xSize*s.ySize; i++ {
			s.cellArray[i].SetValue(s.savedBoard[i])
		}

		s.CheckNeighbors()
		s.FindSafeMoves()
	}

	// If no safe moves are identified, find least risky moves
	if s.safeMoves == 0 {
		arry = s.SelectRiskyMove()
		// fmt.Printf("Debug: Risky move - Column %d Row %d", (arry[0] + 1), (arry[1] + 1))
	} else {
		arry = s.SelectSafeMove()
		s.safeMoves--
	}

	return arry
}

// CheckNeighbors - To determine safe moves, there is a need to examine the
// landscape and determine where the known mines are, so that we can exhaust
// all safe moves in each pass
func (s *solver) CheckNeighbors() {
	for y := 0; y < s.ySize; y++ {
		for x := 0; x < s.xSize; x++ {
			// Check if empty cell
			if s.cellArray[(s.xSize*y)+x].val == 0 {
				continue
				// Check if covered cell
			} else if s.cellArray[(s.xSize*y)+x].val == 0 {
				continue
			}
			s.cellArray[(s.xSize*y)+x].covNeighbors = 0
			s.cellArray[(s.xSize*y)+x].nearbyMines = 0

			// Count the number of covered neighbors
			for yy := y - 1; yy <= y+1; yy++ {
				if yy < 0 || yy >= s.ySize {
					continue
				}
				for xx := x - 1; xx <= x+1; xx++ {
					if xx < 0 || xx >= s.xSize || (yy == y && xx == x) {
						continue
					} else if s.cellArray[(s.xSize*yy)+xx].val == 9 {
						s.cellArray[(s.xSize*y)+x].covNeighbors++
						s.cellArray[(s.xSize*yy)+xx].weight +=
							s.cellArray[(s.xSize*y)+x].val
					}
				}
			}

			// If the neighbors equal the value, all neighbors are mines
			if s.cellArray[(s.xSize*y)+x].val ==
				s.cellArray[(s.xSize*y)+x].covNeighbors {

			}
		}
	}
}

// FindSafeMoves - will handle more deductive means of finding safe moves
func (s *solver) FindSafeMoves() {

}

// DeduceMine - figure out where more of the mines and safe spots are by
// looking at the neighbors' hints
func (s *solver) DeduceMines(y0, x0, y1, x1, y2, x2 int) {

}

// SelectSafeMove - Picks next guaranteed safe moves
func (s *solver) SelectSafeMove() []int {

	return nil
}

// SelectRiskyMove - Weigh the likelihood of each move and select the lowest
// chance of failure
func (s *solver) SelectRiskyMove() []int {

	return nil
}
