// taken from https://github.com/jneander/tic-tac-toe-go/blob/master/ttt/board.go

package tictactoe

type Board struct {
	spaces []string
}

func (b Board) Blank() string {
	return " "
}

func (b *Board) Spaces() []string {
	dup := make([]string, len(b.spaces))
	copy(dup, b.spaces)
	return dup
}

func (b *Board) Mark(pos int, mark string) {
	if pos >= 0 && pos < len(b.spaces) {
		b.spaces[pos] = mark
	}
}

func (b Board) SpacesWithMark(mark string) []int {
	count, result := 0, make([]int, len(b.Spaces()))
	for i, v := range b.Spaces() {
		if v == mark {
			result[count] = i
			count++
		}
	}
	return result[:count]
}

func (b *Board) Reset() {
	b.setBoard()
}

func (board *Board) IsFull() bool {
	for _, mark := range board.Spaces() {
		if mark == board.Blank() {
			return false
		}
	}
	return true
}

func (board *Board) WinningSetExists() (exists bool) {
	for _, set := range solutions() {
		exists = exists || board.allSpacesMatch(set)
	}
	return
}

func (b *Board) WinningMark() (string, bool) {
	spaces := b.Spaces()
	for _, set := range solutions() {
		if b.allSpacesMatch(set) {
			return spaces[set[0]], true
		}
	}
	return "", false
}

func NewBoard() *Board {
	b := new(Board)
	b.setBoard()
	return b
}

// PRIVATE

func (b *Board) setBoard() {
	b.spaces = make([]string, 9)
	for i := range b.spaces {
		b.spaces[i] = " "
	}
}

func (board *Board) allSpacesMatch(pos []int) bool {
	spaces := board.Spaces()
	mark := spaces[pos[0]]
	result := mark != board.Blank()
	for _, i := range pos {
		result = result && spaces[i] == mark
	}
	return result
}

func solutions() [][]int {
	return [][]int{[]int{0, 1, 2}, []int{3, 4, 5}, []int{6, 7, 8},
		[]int{0, 3, 6}, []int{1, 4, 7}, []int{2, 5, 8},
		[]int{0, 4, 8}, []int{2, 4, 6}}
}
