package tictactoe

type Ivan struct {
	Player  *Player
	minimax *Minimax
}

func NewIvan() *Ivan {
	ivan := new(Ivan)
	ivan.minimax = NewMinimax()
	ivan.minimax.SetMinMaxMarks("O", "X")
	ivan.Player = &Player{
		Id:   0,
		Mark: "X",
		Send: make(chan Event, 1),
	}
	return ivan
}

func (c *Ivan) Move(board Board) int {
	moveScores, _ := c.minimax.ScoreAvailableMoves(&board, c.Player.Mark)
	bestMove, bestScore := -1, -1
	for move, score := range moveScores {
		if score > bestScore {
			bestMove, bestScore = move, score
		}
	}
	return bestMove
}

func (c *Ivan) Run(runner *Runner) {
	for {
		e := <-c.Player.Send
		switch e.eventType() {
		case "gameReady":
			if e.playerId() == c.Player.Id {
				boardCp := *runner.Game.Board()
				runner.MakeMove(Move{
					Place:    c.Move(boardCp),
					PlayerId: c.Player.Id,
				})
			}
		case "gameOver":
			runner.Restart(Restart{c.Player.Id})
		case "move":
			if e.playerId() != c.Player.Id {
				boardCp := *runner.Game.Board()
				runner.MakeMove(Move{
					Place:    c.Move(boardCp),
					PlayerId: c.Player.Id,
				})
			}
		default:

		}
	}
}
