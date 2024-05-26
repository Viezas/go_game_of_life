package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Constants for screen dimensions and cell size
const (
	screenWidth  = 1600                    // Width of the game window in pixels
	screenHeight = 400                     // Height of the game window in pixels
	cellSize     = 10                      // Size of each cell in pixels
	rows         = screenHeight / cellSize // Number of rows of cells
	cols         = screenWidth / cellSize  // Number of columns of cells
)

// Cells is a 2D array representing the state of each cell (alive or dead)
type Cells [rows][cols]bool

// Game struct holds the current state of the game
type Game struct {
	cells      Cells // Current state of the cells
	frameCount int   // Counter to control update speed
	isPaused   bool  // Indicates whether the game is paused
}

// Update function is called every frame to update the game state
func (g *Game) Update() error {
	// Check for play/pause toggle input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= screenWidth-100 && x <= screenWidth-50 && y >= 20 && y <= 60 {
			g.isPaused = !g.isPaused
		}
	}

	if !g.isPaused {
		g.frameCount++
		if g.frameCount >= 5 { // Update cells approximately every 80ms (5 frames at 60fps)
			g.cells = makeNextGeneration(g.cells)
			g.frameCount = 0
		}
	}

	return nil
}

// Draw function is called every frame to render the game state
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)   // Clear the screen with black color
	drawCells(g.cells, screen) // Draw the cells on the screen

	// Draw play/pause button
	drawButton(screen, g.isPaused)
}

// Layout function defines the screen dimensions for the game
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight // Return the screen dimensions
}

// generateCells initializes the cells with random values (alive or dead)
func generateCells() Cells {
	var cells Cells
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		for colIndex := 0; colIndex < cols; colIndex++ {
			if rand.Intn(4) == 0 { // 25% chance to be alive
				cells[rowIndex][colIndex] = true
			}
		}
	}
	return cells
}

// makeNextGeneration calculates the next generation of cells based on current state
func makeNextGeneration(generation Cells) Cells {
	var nextGeneration Cells
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		for colIndex := 0; colIndex < cols; colIndex++ {
			neighborCount := calculateNeighborCount(generation, rowIndex, colIndex)
			alive := generation[rowIndex][colIndex]
			if alive && (neighborCount == 2 || neighborCount == 3) {
				nextGeneration[rowIndex][colIndex] = true
			} else if !alive && neighborCount == 3 {
				nextGeneration[rowIndex][colIndex] = true
			} else {
				nextGeneration[rowIndex][colIndex] = false
			}
		}
	}
	return nextGeneration
}

// calculateNeighborCount returns the number of alive neighbors for a given cell
func calculateNeighborCount(cells Cells, currentRow, currentCol int) int {
	rowStart := max(currentRow-1, 0)
	rowEnd := min(currentRow+1, rows-1)
	colStart := max(currentCol-1, 0)
	colEnd := min(currentCol+1, cols-1)
	neighborCount := 0

	for rowIndex := rowStart; rowIndex <= rowEnd; rowIndex++ {
		for colIndex := colStart; colIndex <= colEnd; colIndex++ {
			if rowIndex == currentRow && colIndex == currentCol {
				continue
			}
			if cells[rowIndex][colIndex] {
				neighborCount++
			}
		}
	}
	return neighborCount
}

// drawCells renders the cells on the screen
func drawCells(cells Cells, screen *ebiten.Image) {
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		for colIndex := 0; colIndex < cols; colIndex++ {
			if cells[rowIndex][colIndex] {
				x := float32(colIndex * cellSize)                                           // X coordinate for drawing
				y := float32(rowIndex * cellSize)                                           // Y coordinate for drawing
				vector.DrawFilledRect(screen, x, y, cellSize, cellSize, color.White, false) // Draw a white cell
			}
		}
	}
}

// drawButton renders the play/pause button on the screen
func drawButton(screen *ebiten.Image, isPaused bool) {
	buttonColor := color.RGBA{0x80, 0x80, 0x80, 0xff} // Gray color for the button
	if isPaused {
		vector.DrawFilledRect(screen, float32(screenWidth-100), 20, 50, 40, buttonColor, false)
		vector.DrawFilledRect(screen, float32(screenWidth-90), 30, 10, 20, color.White, false) // Draw play icon
	} else {
		vector.DrawFilledRect(screen, float32(screenWidth-100), 20, 50, 40, buttonColor, false)
		vector.DrawFilledRect(screen, float32(screenWidth-90), 30, 10, 20, color.White, false) // Draw pause icon
		vector.DrawFilledRect(screen, float32(screenWidth-70), 30, 10, 20, color.White, false)
	}
}

// main initializes the game and starts the game loop
func main() {
	game := &Game{
		cells: generateCells(), // Initialize cells with random values
	}

	ebiten.SetWindowSize(screenWidth, screenHeight) // Set the size of the window
	ebiten.SetWindowTitle("Go - Game of Life")      // Set the title of the window
	if err := ebiten.RunGame(game); err != nil {    // Start the game loop
		log.Fatal(err) // Log any errors that occur
	}
}
