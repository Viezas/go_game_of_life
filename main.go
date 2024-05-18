package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Cells [40][160]bool

func main() {
	screen := initScreen()
	cells := generateCells()

	for {
		// Clear screen
		screen.Clear()
		drawCells(cells, screen)

		// Set time between generations
		time.Sleep(80 * time.Millisecond)
		cells = makeNextGeneration(cells)

		// Update screen
		screen.Show()

		if screen.HasPendingEvent() {
			handleEvent(screen)
		}
	}
}

// Init tcell screen
func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)

	return screen
}

// Handle tcell events
func handleEvent(screen tcell.Screen) {
	// Poll event
	event := screen.PollEvent()

	// Process event
	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			quit(screen)
		}
	}
}

// Finish terminal program
func quit(screen tcell.Screen) {
	screen.Fini()
	os.Exit(0)
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

// Give birth to next generation fom provided generation
func makeNextGeneration(generation Cells) Cells {
	var nextGeneration Cells

	for rowIndex := 0; rowIndex < len(generation); rowIndex++ {
		row := generation[rowIndex]

		for colIndex := 0; colIndex < len(row); colIndex++ {
			neighborCount := calculateNeighborCount(generation, rowIndex, colIndex)
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
	return nextGeneration
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

// Draw cells
func drawCells(cells Cells, screen tcell.Screen) {
	style := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		for colIndex := 0; colIndex < len(cells[rowIndex]); colIndex++ {
			if cells[rowIndex][colIndex] {
				// Draw a colored cell for living cell
				screen.SetContent(colIndex, rowIndex, ' ', nil, style)
			}
		}
	}
}
