package tictactoe

import (
	"bytes"
	"fmt"
	"strings"
)

type Event interface {
	eventType() string
	data() string
	playerId() int
}

type Move struct {
	PlayerId int `json:"playerId"`
	Place    int `json:"place"`
}

type Restart struct {
	PlayerId int `json:"playerId" binding:"required"`
}

type NewPlayer struct {
	player *Player
}

type GameReady struct {
	PlayerId int
}

type GameOver struct {
	PlayerId string
}

// move

func (m Move) eventType() string {
	return "move"
}

func (m Move) data() string {
	return fmt.Sprintf("%d:%d", m.PlayerId, m.Place)
}

func (m Move) playerId() int {
	return m.PlayerId
}

// reset

func (r Restart) eventType() string {
	return "restart"
}

func (r Restart) data() string {
	return fmt.Sprintf("%d", r.PlayerId)
}

func (r Restart) playerId() int {
	return r.PlayerId
}

// new player
func (n NewPlayer) eventType() string {
	return "newPlayer"
}

func (n NewPlayer) data() string {
	return fmt.Sprintf("%d:%s", n.player.Id, n.player.Mark)
}

func (n NewPlayer) playerId() int {
	return n.player.Id
}

// game ready

func (g GameReady) eventType() string {
	return "gameReady"
}

func (g GameReady) data() string {
	return fmt.Sprintf("%d", g.PlayerId)
}

func (g GameReady) playerId() int {
	return g.PlayerId
}

// gameover

func (g GameOver) eventType() string {
	return "gameOver"
}

func (g GameOver) data() string {
	return fmt.Sprintf("%s", g.PlayerId)
}

func (g GameOver) playerId() int {
	return -1
}

func ToMessage(m Event) []byte {
	var data bytes.Buffer

	data.WriteString(fmt.Sprintf("event: %s\n", strings.Replace(m.eventType(), "\n", "", -1)))

	lines := strings.Split(m.data(), "\n")
	for _, line := range lines {
		data.WriteString(fmt.Sprintf("data: %s\n", line))
	}

	data.WriteString("\n")
	return data.Bytes()
}
