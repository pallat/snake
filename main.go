package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 320
	screenHeight = 240
	tileSize     = 5
)

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		snake:    NewSnake(),
		food:     NewFood(),
		gameOver: false,
		ticks:    0,
		speed:    10,
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body        []Point
	Direction   Point
	GrowCounter int
}

func NewSnake() *Snake {
	return &Snake{
		Body: []Point{
			{X: screenWidth / tileSize / 2, Y: screenHeight / tileSize / 2},
		},
		Direction: Point{X: 1, Y: 0},
	}
}

func (s *Snake) Move() {
	newHead := Point{
		X: s.Body[0].X + s.Direction.X,
		Y: s.Body[0].Y + s.Direction.Y,
	}
	s.Body = append([]Point{newHead}, s.Body...)

	if s.GrowCounter > 0 {
		s.GrowCounter--
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

type Food struct {
	Position Point
}

func NewFood() *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(screenWidth / tileSize),
			Y: rand.Intn(screenHeight / tileSize),
		},
	}
}

type Game struct {
	snake         *Snake
	food          *Food
	score         int
	gameOver      bool
	ticks         int
	updateCounter int
	speed         int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.restart()
		}
		return nil
	}

	g.updateCounter++
	if g.updateCounter < g.speed {
		return nil
	}
	g.updateCounter = 0

	// Update the snake's position
	g.snake.Move()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: 1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: 1}
	}

	head := g.snake.Body[0]
	if head.X < 0 || head.Y < 0 || head.X >= screenWidth/tileSize || head.Y >= screenHeight/tileSize {
		g.gameOver = true
		g.speed = 10
	}

	for _, part := range g.snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.gameOver = true
			g.speed = 10
		}
	}

	if head.X == g.food.Position.X && head.Y == g.food.Position.Y {
		g.score++
		g.snake.GrowCounter += 1
		g.food = NewFood()
		g.score++
		g.food = NewFood()
		g.snake.GrowCounter = 1

		// Decrease speed (with a lower limit)
		if g.speed > 2 {
			g.speed--
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw snake
	for _, p := range g.snake.Body {
		ebitenutil.DrawRect(screen, float64(p.X*tileSize), float64(p.Y*tileSize), tileSize, tileSize, color.RGBA{0, 255, 0, 255})
	}

	// Draw food
	ebitenutil.DrawRect(screen, float64(g.food.Position.X*tileSize), float64(g.food.Position.Y*tileSize), tileSize, tileSize, color.RGBA{255, 0, 0, 255})

	// Create a font.Face
	face := basicfont.Face7x13

	// Draw game over text
	if g.gameOver {
		text.Draw(screen, "Game Over", face, screenWidth/2-40, screenHeight/2, color.White)
		text.Draw(screen, "Press 'R' to restart", face, screenWidth/2-60, screenHeight/2+16, color.White)
	}

	// Draw score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, face, 5, screenHeight-5, color.White)
}

func (g *Game) restart() {
	g.snake = NewSnake()
	g.score = 0
	g.gameOver = false
	g.food = NewFood()
	g.speed = 10
}
