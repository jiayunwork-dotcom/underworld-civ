package handlers

import (
	"encoding/json"
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

	state, ok := game.GetGameManager().GetGame(gameID)
	if !ok {
		c.Close()
		return
	}

	hub := ws.GetHub()
	if _, ok := hub.GetGameRoom(gameID); !ok {
		hub.CreateGameRoom(gameID, state)
	}

	client := ws.NewClient(c, gameID, playerID)

	hub.Register(client)

	log.Printf("Player %s connected to game %s", playerID, gameID)

	initialMsg, _ := json.Marshal(map[string]interface{}{
		"type": "game_state",
		"data": state,
	})
	client.Send <- initialMsg

	go client.WritePump()
	client.ReadPump()
}
