package ws

import (
	"encoding/json"
	"log"
	"sync"
	"underworld-civ/internal/game"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Client struct {
	ID       string
	PlayerID string
	GameID   string
	Conn     *websocket.Conn
	Send     chan []byte
}

type Hub struct {
	games     map[string]*GameRoom
	register   chan *Client
	unregister chan *Client
	broadcast  chan *GameMessage
	mu         sync.RWMutex
}

type GameRoom struct {
	GameID  string
	Clients map[string]*Client
	State   *game.GameState
}

type GameMessage struct {
	GameID  string
	PlayerID string
	Type    string
	Data    interface{}
}

var hub = NewHub()

func NewHub() *Hub {
	return &Hub{
		games:     make(map[string]*GameRoom),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *GameMessage, 256),
	}
}

func GetHub() *Hub {
	return hub
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if room, ok := h.games[client.GameID]; ok {
				room.Clients[client.ID] = client
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.games[client.GameID]; ok {
				if _, ok := room.Clients[client.ID]; ok {
					delete(room.Clients, client.ID)
					close(client.Send)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			if room, ok := h.games[message.GameID]; ok {
				msgData, _ := json.Marshal(map[string]interface{}{
					"type": message.Type,
					"data": message.Data,
				})
				for _, client := range room.Clients {
					select {
					case client.Send <- msgData:
					default:
						close(client.Send)
						delete(room.Clients, client.ID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) CreateGameRoom(gameID string, state *game.GameState) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.games[gameID] = &GameRoom{
		GameID:  gameID,
		Clients: make(map[string]*Client),
		State:   state,
	}
}

func (h *Hub) RemoveGameRoom(gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.games[gameID]; ok {
		for _, client := range room.Clients {
			close(client.Send)
		}
		delete(h.games, gameID)
	}
}

func (h *Hub) GetGameRoom(gameID string) (*GameRoom, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.games[gameID]
	return room, ok
}

func (h *Hub) BroadcastToGame(gameID, msgType string, data interface{}) {
	h.broadcast <- &GameMessage{
		GameID: gameID,
		Type:   msgType,
		Data:   data,
	}
}

func (h *Hub) SendToPlayer(gameID, playerID, msgType string, data interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.games[gameID]; ok {
		msgData, _ := json.Marshal(map[string]interface{}{
			"type": msgType,
			"data": data,
		})
		for _, client := range room.Clients {
			if client.PlayerID == playerID {
				select {
				case client.Send <- msgData:
				default:
				}
				break
			}
		}
	}
}

func NewClient(conn *websocket.Conn, gameID, playerID string) *Client {
	return &Client{
		ID:       uuid.New().String(),
		PlayerID: playerID,
		GameID:   gameID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func (c *Client) handleMessage(msg []byte) {
	var message struct {
		Action string                 `json:"action"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(msg, &message); err != nil {
		log.Printf("error parsing message: %v", err)
		return
	}

	switch message.Action {
	case "ping":
		c.Send <- []byte(`{"type":"pong"}`)
	}
}
