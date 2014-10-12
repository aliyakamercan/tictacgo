// taken from https://github.com/jneander/tic-tac-toe-go/blob/master/ttt/game.go
package tictactoe

type Game interface {
	Board() *Board
	IsOver() bool
	IsValidMove(int) bool
	ApplyMove(int, string)
	Winner() (string, bool)
}

type game struct {
	board *Board
}

func NewGame() *game {
	g := new(game)
	g.board = NewBoard()
	return g
}

func (g *game) Board() *Board {
	return g.board
}

func (g *game) IsOver() bool {
	return g.board.WinningSetExists() || g.board.IsFull()
}

func (g *game) Winner() (string, bool) {
	return g.board.WinningMark()
}

func (g *game) IsValidMove(space int) bool {
	board := g.board
	isInRange := space >= 0 && space < len(board.Spaces())
	return isInRange && board.Spaces()[space] == board.Blank()
}

func (g *game) ApplyMove(pos int, mark string) {
	if g.IsValidMove(pos) {
		g.board.Mark(pos, mark)
	}
}

func (g *game) Reset() {
	g.board.Reset()
}
