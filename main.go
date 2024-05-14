package main

import "fmt"

func main() {
	firstGeneration := [3][3]bool{{true, false, true}, {false, true, false}, {false, false, true}}
	var nextGeneration [3][3]bool

	for rowIndex := 0; rowIndex < len(firstGeneration); rowIndex++ {
		row := firstGeneration[rowIndex]

		for colIndex := 0; colIndex < len(row); colIndex++ {
			neighborCount := calculateNeighborCount()
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

func calculateNeighborCount() int {
	return 5
}

func printCells(cells [3][3]bool) {
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		row := cells[rowIndex]
		for colIndex := 0; colIndex < len(row); colIndex++ {
			if row[colIndex] {
				fmt.Print("O")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}
