package main

import (
	"fmt"
	"os"
	"time"

	"github.com/foxyblue/gogue/gogue/area"
	"github.com/foxyblue/gogue/gogue/creature"
	"github.com/foxyblue/gogue/gogue/feed"
	"github.com/gdamore/tcell"
)

// Game holds the instance of the game
type Game struct {
	Screen tcell.Screen

	// Level refers to the level at which the active area exists
	Level int

	// ActiveArea refers to the active area to which the player is in.
	ActiveArea *area.Area

	// Player refers to the user
	Player *creature.Player

	// stdFeed is the in game feed
	StdFeed *feed.Feed
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		Screen: newScreen(),
		// Player: newPlayer(),
		Level: 0,
	}
}

// newScreen generates a screen instance for the game to be played
func newScreen() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	return s
}

// CreateArea creates a new playable area
func (game *Game) CreateArea() {
	area := area.NewArea(game.Level, game.Screen)
	game.ActiveArea = area
}

// CreatePlayer creates a user
func (game *Game) CreatePlayer(x, y int) {
	player := creature.NewPlayer(x, y)
	game.Player = player
}

// Draw will render the game on screen
func (game *Game) Draw() {
	game.Screen.Clear()
	game.ActiveArea.Draw()
	st := tcell.StyleDefault
	p := game.Player.Creature
	game.Screen.SetCell(p.X, p.Y, st.Background(p.Color), p.Appearance)
	game.StdFeed.Draw()
	game.Screen.Show()
}

// CreateFeed creates a new message queue
func (game *Game) CreateFeed() {
	f := feed.NewFeed(game.Screen)
	game.StdFeed = f
}

func main() {
	game := NewGame()
	game.CreateArea()
	x := game.ActiveArea.Start.X
	y := game.ActiveArea.Start.Y
	game.CreatePlayer(x, y)
	game.CreateFeed()
	game.StdFeed.Feed.Enqueue("This is text")

	// This is the Key Listener Channel
	quit := make(chan struct{})
	go func() {
		for {
			ev := game.Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					game.Screen.Sync()
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'k':
						game.Player.Creature.Move(0, -1)
					case 'j':
						game.Player.Creature.Move(0, 1)
					case 'h':
						game.Player.Creature.Move(-1, 0)
					case 'l':
						game.Player.Creature.Move(1, 0)
					}
					game.Draw()
				}
			case *tcell.EventResize:
				game.Screen.Sync()
			}
		}
	}()

	// Main Gameloop
	cnt := 0
	dur := time.Duration(0)
gameLoop:
	for {
		select {
		case <-quit:
			break gameLoop
		case <-time.After(time.Millisecond * 50):
		}
		start := time.Now()
		game.Draw()
		cnt++
		dur += time.Now().Sub(start)
	}

	game.Screen.Fini()
	fmt.Println("Game has ended.")
}
