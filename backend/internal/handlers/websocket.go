package handlers

import (
	"log"
	"underworld-civ/internal/game"
	"underworld-civ/internal/ws"

	"github.com/gofiber/contrib/websocket"
)

func WebSocketHandler(c *websocket.Conn) {
	gameID := c.Params("id")
	playerID := c.Query("player_id")

	if gameID == "" || playerID == "" {
		c.Close()
		return
	}

	_, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		c.Close()
		return
	}

	client := ws.NewClient(c, gameID, playerID)
	hub := ws.GetHub()

	hub.Register(client)

	log.Printf("Player %s connected to game %s", playerID, gameID)

	go client.WritePump()
	client.ReadPump()
}
