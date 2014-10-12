package tictactoe

import (
	"math/rand"
)

type Runner struct {
	Game Game `json:"game"`

	reset int

	// registered players
	players map[*Player]bool

	// inbound
	broadcast chan Event

	// register players
	register chan *Player

	// unregister players
	unregister chan *Player
}

func NewRunner() *Runner {
	runner := &Runner{
		broadcast:  make(chan Event, 1),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		players:    make(map[*Player]bool),
	}
	return runner
}

func (runner *Runner) IsFull() bool {
	if len(runner.players) == 2 {
		return true
	}
	return false
}

func (runner *Runner) NewPlayer() *Player {
	id := len(runner.players)
	mark := "X"
	if id == 1 {
		mark = "O"
		defer runner.NewRound()
	}
	player := &Player{
		Id:   id,
		Mark: mark,
		Send: make(chan Event, 1),
	}
	runner.register <- player
	runner.broadcast <- &NewPlayer{player}
	return player
}

func (runner *Runner) AgainstIvan() *Player {
	ivan := NewIvan()
	runner.register <- ivan.Player
	go ivan.Run(runner)
	runner.broadcast <- &NewPlayer{ivan.Player}
	return ivan.Player
}

func (runner *Runner) Leave(p *Player) int {
	runner.unregister <- p
	// broadcast player left
	return len(runner.players)
}

func (runner *Runner) NewRound() {
	runner.Game = NewGame()
	readyEvent := &GameReady{
		PlayerId: rand.Intn(1),
	}
	runner.broadcast <- readyEvent
}

func (h *Runner) Run() {
	for {
		select {
		case c := <-h.register:
			h.players[c] = true
		case c := <-h.unregister:
			if _, ok := h.players[c]; ok {
				delete(h.players, c)
				close(c.Send)
			}
		case m := <-h.broadcast:
			for c := range h.players {
				select {
				case c.Send <- m:
				default:
					delete(h.players, c)
					close(c.Send)
				}
			}
		}
	}
}

func (runner *Runner) MakeMove(move Move) {
	var current *Player
	for p := range runner.players {
		if p.Id == move.PlayerId {
			current = p
			break
		}
	}
	if runner.Game.IsValidMove(move.Place) && !runner.Game.IsOver() {
		runner.Game.ApplyMove(move.Place, current.Mark)
		runner.broadcast <- move
		if runner.Game.IsOver() {
			winner, _ := runner.Game.Winner()
			runner.broadcast <- &GameOver{winner}
		}
	}

}

func (runner *Runner) Restart(r Restart) {
	runner.broadcast <- r
	runner.reset++
	if runner.reset == 2 {
		runner.reset = 0
		runner.NewRound()
	}
}
