package main

import (
	"fmt"
	"math/rand"
)

type Cells [20][80]bool

func main() {
	firstGeneration := generateCells()
	var nextGeneration Cells

	for rowIndex := 0; rowIndex < len(firstGeneration); rowIndex++ {
		row := firstGeneration[rowIndex]

		for colIndex := 0; colIndex < len(row); colIndex++ {
			neighborCount := calculateNeighborCount(firstGeneration, rowIndex, colIndex)
			alive := row[colIndex]

			if alive && (neighborCount == 2 || neighborCount == 3) {
				// KEEP CELL ALIVE
				nextGeneration[rowIndex][colIndex] = true
			} else if !alive && neighborCount == 3 {
				// REVIVE DEAD CELL
				nextGeneration[rowIndex][colIndex] = true
			} else {
				// KILL CELL because of LONELINESS or OVERPOPULATION or cell was ALREADY DEAD
				nextGeneration[rowIndex][colIndex] = false
			}
		}
	}

	fmt.Println("First generation:")
	printCells(firstGeneration)

	fmt.Println("Next generation:")
	printCells(nextGeneration)
}

// Generate cells with random values
func generateCells() Cells {
	// Initialize cells with default false value
	var cells Cells

	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		for colIndex := 0; colIndex < len(cells[rowIndex]); colIndex++ {
			// Generate living cell with 25% chance
			if rand.Intn(4) == 0 {
				cells[rowIndex][colIndex] = true
			}
		}
	}

	return cells
}

// Calculate and return the number of neighbors for a given cell
func calculateNeighborCount(cells Cells, currentRow, currentCol int) int {
	rowStart := max(currentRow-1, 0)
	rowEnd := min(currentRow+1, len(cells)-1)
	colStart := max(currentCol-1, 0)
	colEnd := min(currentCol+1, len(cells[0])-1)
	neighborCount := 0

	for rowIndex := rowStart; rowIndex <= rowEnd; rowIndex++ {
		for colIndex := colStart; colIndex <= colEnd; colIndex++ {
			isRefCell := rowIndex == currentRow && colIndex == currentCol

			// Increase neighbor count if this is not our reference cell and there is a living neighbor.
			if !isRefCell && cells[rowIndex][colIndex] {
				neighborCount++
			}
		}
	}

	return neighborCount
}

// Print readable cells
func printCells(cells Cells) {
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		row := cells[rowIndex]
		for colIndex := 0; colIndex < len(row); colIndex++ {
			if row[colIndex] {
				fmt.Print("*")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
