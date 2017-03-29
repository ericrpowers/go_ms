package main

import (
	"log"
	"math/rand"
	"time"
)

// The value representing a mine in the field
const mineVal = -1

// Minefield - struct
type Minefield struct {
	minefield []int
}

// NewMineField - constructor
func NewMineField(colLength, rowLength, mineCount, row, column int) *Minefield {
	m := new(Minefield)
	m.minefield = make([]int, colLength*rowLength)
	m.placeMines(colLength, rowLength, mineCount, row, column)
	m.fillHints(colLength, rowLength)

	return m
}

// GetMineVal - Gets int representation of a mine in the minefield
func (m *Minefield) GetMineVal() int {
	return mineVal
}

// GetCellVal - Retrieves value of the specified cell
func (m *Minefield) GetCellVal(pos int) int {
	return m.minefield[pos]
}

// Avoid entries with mines and the first selected position
func (m *Minefield) placeMines(colLength, rowLength, mineCount, row,
	column int) {
	if colLength*rowLength == mineCount {
		log.Fatal("Number of mines equal size of board")
	}
	var r, c int
	for i := 0; i < mineCount; i++ {
		rGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		r, c = rGen.Intn(rowLength), rGen.Intn(colLength)
		for m.minefield[(colLength*r)+c] == mineVal ||
			(r == row && c == column) {
			r = rGen.Intn(rowLength)
			c = rGen.Intn(colLength)
		}

		m.minefield[(colLength*r)+c] = mineVal
	}
}

/*  To fill in the hints, need to look at every empty cell and count the Number
    of adjacent mines into that cell. Make sure to avoid going outside the
    boundaries of the playing field */
func (m *Minefield) fillHints(colLength, rowLength int) {
	for row := 0; row < rowLength; row++ {
		for column := 0; column < colLength; column++ {
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if row+i >= 0 && row+i < rowLength && column+j >= 0 &&
						column+j < colLength {
						if m.minefield[(colLength*row)+column] != mineVal {
							if m.minefield[(colLength*(row+i))+(column+j)] == mineVal {
								m.minefield[(colLength*row)+column]++
							}
						}
					}
				}
			}
		}
	}
}
