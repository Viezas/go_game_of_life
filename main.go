package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Cells [][]bool

func main() {
	screen := initScreen()
	width, height := 40, 20 // Taille plus petite de la grille initiale
	cells := generateCells(width, height)
	speed := 1

	for {
		// Clear screen
		screen.Clear()
		drawCells(cells, screen)

		// Set time between generations
		time.Sleep(time.Duration(80/speed) * time.Millisecond)
		cells = makeNextGeneration(cells)

		// Update screen
		screen.Show()

		if screen.HasPendingEvent() {
			width, height, cells, speed = handleEvent(screen, width, height, cells, speed)
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
func handleEvent(screen tcell.Screen, width, height int, cells Cells, speed int) (int, int, Cells, int) {
	// Poll event
	event := screen.PollEvent()

	// Process event
	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Sync()
		newWidth, newHeight := screen.Size()
		cells = resizeCells(cells, newWidth, newHeight)
		width, height = newWidth, newHeight
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			quit(screen)
		}
		// Adjust speed
		switch event.Rune() {
		case '1':
			speed = 1
		case '2':
			speed = 2
		case '4':
			speed = 4
		}
	}
	return width, height, cells, speed
}

// Finish terminal program
func quit(screen tcell.Screen) {
	screen.Fini()
	os.Exit(0)
}

// Generate cells with random values
func generateCells(width, height int) Cells {
	// Initialize cells with default false value
	cells := make(Cells, height)
	for rowIndex := 0; rowIndex < height; rowIndex++ {
		cells[rowIndex] = make([]bool, width)
		for colIndex := 0; colIndex < width; colIndex++ {
			// Generate living cell with 25% chance
			if rand.Intn(4) == 0 {
				cells[rowIndex][colIndex] = true
			}
		}
	}

	return cells
}

// Resize cells to match new dimensions
func resizeCells(oldCells Cells, width, height int) Cells {
	newCells := make(Cells, height)
	for rowIndex := 0; rowIndex < height; rowIndex++ {
		newCells[rowIndex] = make([]bool, width)
		for colIndex := 0; colIndex < width; colIndex++ {
			if rowIndex < len(oldCells) && colIndex < len(oldCells[0]) {
				newCells[rowIndex][colIndex] = oldCells[rowIndex][colIndex]
			}
		}
	}
	return newCells
}

// Give birth to next generation fom provided generation
func makeNextGeneration(generation Cells) Cells {
	height := len(generation)
	width := len(generation[0])
	var nextGeneration Cells
	nextGeneration = make(Cells, height)
	for i := range nextGeneration {
		nextGeneration[i] = make([]bool, width)
	}

	for rowIndex := 0; rowIndex < height; rowIndex++ {
		row := generation[rowIndex]

		for colIndex := 0; colIndex < width; colIndex++ {
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

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
