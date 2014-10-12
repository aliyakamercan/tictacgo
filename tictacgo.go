package main

import (
	"github.com/aliyakamercan/tictacgo/tictactoe"
	"github.com/gin-gonic/gin"
	"math/rand"
	//	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	runners := make(map[string]*tictactoe.Runner)

	r := gin.Default()

	r.GET("/game/:runnerid/updates", func(c *gin.Context) {
		runnerid := c.Params.ByName("runnerid")
		runner, ok := runners[runnerid]

		if !ok {
			c.String(404, "Not Found")
			return
		}

		// f, fok := c.Writer.(http.Flusher)
		// if !fok {
		// 	http.Error(c.Writer, "streaming unsupported", http.StatusInternalServerError)
		// 	return
		// }

		// game updates (server sent events)
		c.Writer.Header().Add("Content-Type", "text/event-stream")
		c.Writer.Header().Add("Cache-Control", "no-cache")
		c.Writer.Header().Add("Connection", "keep-alive")
		c.Writer.WriteHeader(200)

		if !runner.IsFull() {
			player := runner.NewPlayer()
			notify := c.Writer.CloseNotify()
			for {
				select {
				case event := <-player.Send:
					println(player)
					println(string(tictactoe.ToMessage(event)))
					c.Writer.Write(tictactoe.ToMessage(event))
					c.Writer.Flush()
				case <-notify:
					if runner.Leave(player) == 0 {
						delete(runners, runnerid)
					}
					return
				}
			}
		} else {
			c.String(400, "This game is full")
			return
		}
	})

	r.GET("/game/:runnerid", func(c *gin.Context) {
		runnerid := c.Params.ByName("runnerid")
		runner, ok := runners[runnerid]
		if !ok {
			c.String(404, "Not Found")
			return
		}
		c.JSON(200, runner)
	})

	r.PATCH("/game/:gameid/move", func(c *gin.Context) {
		// handle moves
		var move tictactoe.Move
		c.Bind(&move)

		runnerid := c.Params.ByName("gameid")
		runner := runners[runnerid]

		runner.MakeMove(move)

		c.JSON(200, gin.H{"status": "success"})
	})

	r.PATCH("/game/:gameid/restart", func(c *gin.Context) {
		// handle moves
		var r tictactoe.Restart
		c.Bind(&r)

		runnerid := c.Params.ByName("gameid")
		runner := runners[runnerid]

		runner.Restart(r)

		c.JSON(200, gin.H{"status": "success"})
	})

	r.POST("/game/", func(c *gin.Context) {
		runner := tictactoe.NewRunner()
		runnerid := randSeq(6)
		runners[runnerid] = runner
		go runner.Run()
		c.JSON(200, gin.H{"id": runnerid})
	})

	r.POST("/game/ivan", func(c *gin.Context) {
		runner := tictactoe.NewRunner()
		runnerid := randSeq(6)
		runners[runnerid] = runner
		go runner.Run()
		runner.AgainstIvan()
		c.JSON(200, gin.H{"id": runnerid})
	})

	r.GET("/", func(c *gin.Context) {
		c.File("public/index.html")
	})

	r.Static("/static", "public")

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
