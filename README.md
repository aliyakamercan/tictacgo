TicTacGo
=========

It is a tictactoe written in go using [ServerSentEvent](https://developer.mozilla.org/en-US/docs/Server-sent_events/Using_server-sent_events).
  - can be played against an unbeatle computer
  - can be played againsts other people

You can find a live demo on [here](http://protected-dawn-5794.herokuapp.com/)

How To Run
---------

** Requirements**
   - Go 1.3 (Download from golang.org/dl & set you GOPATH)

``` 
go get github.com/aliyakamercan/tictacgo
cd $GOPATH/src/github.com/aliyakamercan/tictacgo
go run web.go
```

Todo
---

* Tests
* Player left event
* Redesign Events and some other parts (this was my first go proj)

Kudos
---

- I ripped off almost all styles from http://perfecttictactoe.herokuapp.com/
- I reused some parts from [jneander/tic-tac-toe-go](https://github.com/jneander/tic-tac-toe-go)
