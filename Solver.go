package main

import (
	"fmt"
	"log"
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
			fmt.Println(coordArray)
		}

		if s.Win() {
			s.wins++
		} else {
			s.ShowBoard()
			log.Fatal("Debug: Lost game")
		}
		iterations--
	}
	elapsedTime := time.Since(startTime).Seconds()

	// Now to print out the statistics
	fmt.Printf("%20s %10d\n", "Total games", s.numOfGames)
	fmt.Printf("%20s %10d\n", "Total wins", s.wins)
	fmt.Printf("%20s %10f\n", "Total time taken (s)", elapsedTime)
	fmt.Printf("%20s %10f\n", "Win Percentage (%)", float64(s.wins)/(float64(s.numOfGames))*100)
	fmt.Printf("%20s %10f\n", "Average # of turns", float64(s.numOfTurns)/float64(s.numOfGames))
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
		fmt.Printf("Debug: Risky move - Column %d Row %d", (arry[0] + 1), (arry[1] + 1))
	} else {
		arry = s.SelectSafeMove()
		s.safeMoves--
	}

	return arry
}

func validateCell(i, iSize int) bool {
	if i < 0 || i >= iSize {
		return false
	}
	return true
}

func validateCells(size, x, xx, y, yy int) bool {
	if xx < 0 || xx >= size || (yy == y && xx == x) {
		return false
	}
	return true
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
				if !validateCell(yy, s.ySize) {
					continue
				}
				for xx := x - 1; xx <= x+1; xx++ {
					if !validateCells(s.xSize, x, xx, y, yy) {
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
				for yy := y - 1; yy <= y+1; yy++ {
					if !validateCell(yy, s.ySize) {
						continue
					}
					for xx := x - 1; xx <= x+1; xx++ {
						if !validateCells(s.xSize, x, xx, y, yy) {
							continue
						} else if s.cellArray[(s.xSize*yy)+xx].isMine == 1 {
							s.cellArray[(s.xSize*y)+x].nearbyMines++
						}
					}
				}
			}

			// See if we can already identify safe moves
			if s.cellArray[(s.xSize*y)+x].val ==
				s.cellArray[(s.xSize*y)+x].nearbyMines {
				for yy := y - 1; yy <= y+1; yy++ {
					if !validateCell(yy, s.ySize) {
						continue
					}
					for xx := x - 1; xx <= x+1; xx++ {
						if !validateCells(s.xSize, x, xx, y, yy) {
							continue
						} else if s.cellArray[(s.xSize*yy)+xx].val == 9 &&
							s.cellArray[(s.xSize*yy)+xx].isMine == -1 {
							s.cellArray[(s.xSize*yy)+xx].isMine = 0
							s.safeMoves++
						}
					}
				}
			}

			// If more mines then value, we ran into an issue
			if s.cellArray[(s.xSize*y)+x].val < s.cellArray[(s.xSize*y)+x].nearbyMines {
				s.ShowBoard()
				log.Fatalf("ERROR - more mines (%d) than value (%d)! Row %d Colume %d\n",
					s.cellArray[(s.xSize*y)+x].nearbyMines, s.cellArray[(s.xSize*y)+x].val,
					y+1, x+1)
			}
		}
	}
}

// FindSafeMoves - will handle more deductive means of finding safe moves
func (s *solver) FindSafeMoves() {
	for y := 0; y < s.ySize; y++ {
		for x := 0; x < s.xSize; x++ {
			x1, x2, y1, y2 := -1, -1, -1, -1
			if s.cellArray[(s.xSize*y)+x].val == 0 || s.cellArray[(s.xSize*y)+x].val == 9 {
				continue
			}

			/* If we have a scenario where we know val - 1 mines, we should be able to
			determine the other mines and safe moves based on neighboring hints */
			if s.cellArray[(s.xSize*y)+x].val-s.cellArray[(s.xSize*y)+x].nearbyMines == 1 &&
				s.cellArray[(s.xSize*y)+x].covNeighbors-s.cellArray[(s.xSize*y)+x].nearbyMines == 2 {
				for yy := 0; yy <= y+1; yy++ {
					if !validateCell(yy, s.ySize) {
						continue
					}
					for xx := 0; xx < x+1; xx++ {
						if !validateCells(s.xSize, x, xx, y, yy) {
							continue
						} else if s.cellArray[(s.xSize*yy)+xx].val == 9 &&
							s.cellArray[(s.xSize*yy)+xx].isMine == -1 {
							if x1 == -1 {
								x1, y1 = xx, yy
							} else {
								x2, y2 = xx, yy
							}
						}
					}
				}

				if x1 == -1 || x2 == -1 {
					continue
				}
				// Look below and above cell
				if x1 == x2 {
					if y-1 >= 0 && y-1 < s.ySize {
						s.DeduceMines(y-1, x, y1, x1, y2, x2)
					}
					if y+1 >= 0 && y+1 < s.ySize {
						s.DeduceMines(y+1, x, y1, x1, y2, x2)
					}
				}
				// Look left and right of cell
				if y1 == y2 {
					if x-1 >= 0 && x-1 < s.xSize {
						s.DeduceMines(y, x-1, y1, x1, y2, x2)
					}
					if x+1 >= 0 && x+1 < s.xSize {
						s.DeduceMines(y, x+1, y1, x1, y2, x2)
					}
				}
			}
		}
	}
}

// DeduceMine - figure out where more of the mines and safe spots are by
// looking at the neighbors' hints
func (s *solver) DeduceMines(y0, x0, y1, x1, y2, x2 int) {
	for yy := y0 - 1; yy <= y0+1; yy++ {
		if !validateCell(yy, s.ySize) {
			continue
		}
		for xx := x0 - 1; xx <= x0+1; xx++ {
			// Ignore entries we have already looked at or already identified as mines
			if !validateCells(s.xSize, x0, xx, y0, yy) || (yy == y1 && xx == x1) ||
				(yy == y2 && xx == x2) || !(s.cellArray[(s.xSize*yy)+xx].val == 9 &&
				s.cellArray[(s.xSize*yy)+xx].isMine == -1) {
				continue
			}
			if s.cellArray[(s.xSize*y0)+x0].val-s.cellArray[(s.xSize*yy)+xx].nearbyMines == 1 {
				// Return if any of the coordinates are 3 away from the cells in question
				if yy == y1+3 || yy == y1-3 || yy == y2+3 || yy == y2-3 ||
					xx == x1+3 || xx == x1-3 || xx == x2+3 || xx == x2-3 {
					return
				}
				s.cellArray[(s.xSize*yy)+xx].isMine = 0
				s.safeMoves++
				if s.cellArray[(s.xSize*y0)+x0].covNeighbors-s.cellArray[(s.xSize*y0)+x0].val == 2 &&
					s.cellArray[(s.xSize*yy)+xx].isMine == -1 {
					s.cellArray[(s.xSize*yy)+xx].isMine = 1

					// Increase the nearby mine count if val is between 1 and 8
					for yyy := yy - 1; yyy <= yy+1; yyy++ {
						if !validateCell(yyy, s.ySize) {
							continue
						}
						for xxx := xx - 1; xxx <= xx+1; xxx++ {
							if !validateCells(s.xSize, xx, xxx, yy, yyy) {
								continue
							}
							if s.cellArray[(s.xSize*yyy)+xxx].val > 0 || s.cellArray[(s.xSize*yyy)+xxx].val < 9 {
								s.cellArray[(s.xSize*yyy)+xxx].nearbyMines++
							}
						}
					}
				}
			}
		}
	}
}

// SelectSafeMove - Picks next guaranteed safe moves
func (s *solver) SelectSafeMove() []int {
	// All safe moves are identified with isMine == 0
	for y := 0; y < s.ySize; y++ {
		for x := 0; x < s.xSize; x++ {
			if s.cellArray[(s.xSize*y)+x].val != 9 ||
				s.cellArray[(s.xSize*y)+x].isMine != 0 {
				continue
			}
			s.cellArray[(s.xSize*y)+x].isMine = -1
			return []int{x, y}
		}
	}
	log.Fatal("ERROR - Could not find a safe move!")

	return nil
}

// SelectRiskyMove - Weigh the likelihood of each move and select the lowest
// chance of failure
func (s *solver) SelectRiskyMove() []int {
	prob := 99
	var array []int      // least risky cell
	var blindArray []int // blind click cell
	for y := 0; y < s.ySize; y++ {
		for x := 0; x < s.xSize; x++ {
			if s.cellArray[(s.xSize*y)+x].val != 9 ||
				s.cellArray[(s.xSize*y)+x].isMine != -1 {
				continue
			}
			cWeight := s.cellArray[(s.xSize*y)+x].weight
			if cWeight > 0 && cWeight < prob {
				prob = cWeight
				array = []int{x, y}
			} else if cWeight <= 0 {
				blindArray = []int{x, y}
			}
		}
	}

	if array == nil && blindArray == nil {
		s.ShowBoard()
		log.Fatal("ERROR - Could not find a risky move!")
	} else if array == nil {
		// Fall back on blind clicking as it is the last option
		array = blindArray
	}

	return array
}
