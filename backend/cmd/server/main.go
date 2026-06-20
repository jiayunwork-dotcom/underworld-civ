package main

import (
	"log"
	"time"
	"underworld-civ/internal/cache"
	"underworld-civ/internal/config"
	"underworld-civ/internal/db"
	"underworld-civ/internal/game"
	"underworld-civ/internal/handlers"
	"underworld-civ/internal/ws"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	if err := db.Init(cfg); err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		log.Println("Database connected successfully")
	}

	if err := cache.Init(cfg); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	go ws.GetHub().Run()

	go gameStateBroadcaster()

	app := fiber.New(fiber.Config{
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(logger.New())

	api := app.Group("/api")

	api.Get("/games", handlers.ListGames)
	api.Post("/games", handlers.CreateGame)
	api.Get("/games/:id", handlers.GetGame)
	api.Post("/games/:id/join", handlers.JoinGame)
	api.Post("/games/:id/start", handlers.StartGame)
	api.Post("/games/:id/actions", handlers.SubmitAction)

	api.Get("/races", handlers.GetRaces)
	api.Get("/buildings", handlers.GetBuildings)
	api.Get("/units", handlers.GetUnits)
	api.Get("/techs", handlers.GetTechs)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/games/:id", websocket.New(handlers.WebSocketHandler))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

func gameStateBroadcaster() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		hub := ws.GetHub()
		for _, room := range getGameRooms() {
			if state, ok := game.GetGameManager().GetGame(room); ok {
				hub.BroadcastToGame(room, "game_state", state)
			}
		}
	}
}

func getGameRooms() []string {
	hub := ws.GetHub()
	rooms := make([]string, 0)

	games := game.GetGameManager().ListGames()
	for _, g := range games {
		if _, ok := hub.GetGameRoom(g.ID); ok {
			rooms = append(rooms, g.ID)
		}
	}

	return rooms
}
