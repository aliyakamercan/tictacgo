package tictactoe

type Player struct {
	Id   int    `json:"id"`
	Mark string `json:"mark"`
	Send chan Event
}
