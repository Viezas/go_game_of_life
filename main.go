package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Constants for cell size
const (
	cellSize = 10 // Size of each cell in pixels
)

// Cells is a 2D array representing the state of each cell (alive or dead)
type Cells [][]bool

// Game struct holds the current state of the game
type Game struct {
	cells      Cells // Current state of the cells
	frameCount int   // Counter to control update speed
	isPaused   bool  // Indicates whether the game is paused
	speed      int   // Speed multiplier for the game update
	width      int   // Current width of the window
	height     int   // Current height of the window
	rows       int   // Number of rows of cells
	cols       int   // Number of columns of cells
	drawing    bool  // Indicates whether we are in drawing mode
}

// Update function is called every frame to update the game state
func (g *Game) Update() error {
	// Check for play/pause toggle input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= g.width-100 && x <= g.width-50 && y >= 20 && y <= 60 {
			g.isPaused = !g.isPaused
		} else if g.drawing {
			row := y / cellSize
			col := x / cellSize
			if row >= 0 && row < g.rows && col >= 0 && col < g.cols {
				g.cells[row][col] = !g.cells[row][col]
			}
		}
	}

	// Adjust speed with number keys
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.speed = 1
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.speed = 2
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.speed = 4
	}

	// Toggle drawing mode with the D key
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.drawing = !g.drawing
	}

	// Load a pattern from file with the L key
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.loadPatternFromFile("pattern.json")
	}

	if !g.isPaused && !g.drawing {
		g.frameCount++
		if g.frameCount >= 5/g.speed { // Update cells based on the speed multiplier
			g.cells = makeNextGeneration(g.cells, g.rows, g.cols)
			g.frameCount = 0
		}
	}

	return nil
}

// Draw function is called every frame to render the game state
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)                   // Clear the screen with black color
	drawCells(g.cells, screen, g.rows, g.cols) // Draw the cells on the screen

	// Draw play/pause button
	drawButton(screen, g.isPaused, g.drawing, g.width)
}

// Layout function defines the screen dimensions for the game
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	g.rows = outsideHeight / cellSize
	g.cols = outsideWidth / cellSize
	g.cells = resizeCells(g.cells, g.rows, g.cols)
	return outsideWidth, outsideHeight // Return the screen dimensions
}

// generateCells initializes the cells with random values (alive or dead)
func generateCells(rows, cols int) Cells {
	cells := make(Cells, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		cells[rowIndex] = make([]bool, cols)
		for colIndex := 0; colIndex < cols; colIndex++ {
			if rand.Intn(4) == 0 { // 25% chance to be alive
				cells[rowIndex][colIndex] = true
			}
		}
	}
	return cells
}

// resizeCells resizes the cells array to the new dimensions, preserving existing cells and generating new ones if needed
func resizeCells(cells Cells, rows, cols int) Cells {
	newCells := make(Cells, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		newCells[rowIndex] = make([]bool, cols)
		for colIndex := 0; colIndex < cols; colIndex++ {
			if rowIndex < len(cells) && colIndex < len(cells[rowIndex]) {
				newCells[rowIndex][colIndex] = cells[rowIndex][colIndex]
			} else {
				// Generate new cells with a 25% chance to be alive
				newCells[rowIndex][colIndex] = rand.Intn(4) == 0
			}
		}
	}
	return newCells
}

// makeNextGeneration calculates the next generation of cells based on current state
func makeNextGeneration(generation Cells, rows, cols int) Cells {
	nextGeneration := make(Cells, rows)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		nextGeneration[rowIndex] = make([]bool, cols)
		for colIndex := 0; colIndex < cols; colIndex++ {
			neighborCount := calculateNeighborCount(generation, rowIndex, colIndex, rows, cols)
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
func calculateNeighborCount(cells Cells, currentRow, currentCol, rows, cols int) int {
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
func drawCells(cells Cells, screen *ebiten.Image, rows, cols int) {
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
func drawButton(screen *ebiten.Image, isPaused, isDrawing bool, width int) {
	buttonColor := color.RGBA{0x80, 0x80, 0x80, 0xff} // Gray color for the button
	vector.DrawFilledRect(screen, float32(width-100), 20, 50, 40, buttonColor, false)
	if isPaused {
		if isDrawing {
			// Draw drawing mode indicator
			vector.DrawFilledRect(screen, float32(width-90), 30, 30, 20, color.RGBA{0x00, 0xFF, 0x00, 0xFF}, false) // Draw green rectangle for drawing mode
		} else {
			// Draw play icon
			vector.DrawFilledRect(screen, float32(width-90), 30, 10, 20, color.White, false) // Left part of play icon
		}
	} else {
		// Draw pause icon
		vector.DrawFilledRect(screen, float32(width-90), 30, 10, 20, color.White, false) // Left part of pause icon
		vector.DrawFilledRect(screen, float32(width-70), 30, 10, 20, color.White, false) // Right part of pause icon
	}
}

// loadPatternFromFile loads a pattern from a JSON file
func (g *Game) loadPatternFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return
	}

	var pattern Cells
	err = json.Unmarshal(data, &pattern)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return
	}

	// Resize cells if necessary
	g.cells = resizeCells(pattern, g.rows, g.cols)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// main initializes the game and starts the game loop
func main() {
	game := &Game{
		cells: generateCells(40, 80), // Initialize cells with random values
		speed: 1,                      // Initial speed multiplier
	}

	ebiten.SetWindowSize(800, 400)                                 // Set the initial size of the window
	ebiten.SetWindowTitle("Go - Game of Life")                     // Set the title of the window
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled) // Allow window to be resizable
	if err := ebiten.RunGame(game); err != nil {                   // Start the game loop
		log.Fatal(err) // Log any errors that occur
	}
}
