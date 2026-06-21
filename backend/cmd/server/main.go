package main

import (
	"log"
	"time"
	"underworld-civ/internal/cache"
	"underworld-civ/internal/config"
	"underworld-civ/internal/db"
	"underworld-civ/internal/game"
	"underworld-civ/internal/handlers"
	"underworld-civ/internal/models"
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

	setupGameManagerCallbacks()

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
	api.Get("/games/:id/tech-tree", handlers.GetTechTree)
	api.Post("/games/:id/tech-tree", handlers.SetResearch)

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
		games := game.GetGameManager().ListGames()
		for _, g := range games {
			if g.Status != "waiting" {
				if _, ok := hub.GetGameRoom(g.ID); !ok {
					hub.CreateGameRoom(g.ID, g)
				}
				hub.BroadcastToGame(g.ID, "game_state", g)
			}
		}
	}
}

func setupGameManagerCallbacks() {
	gm := game.GetGameManager()

	gm.SyncTechs = func(gameID string, playerID string, techs []string, currentResearch string, username string) error {
		return cache.SyncPlayerTechs(gameID, playerID, techs, currentResearch, username)
	}

	gm.SaveGameRecord = func(gameID string, turnNumber int, gameStateJSON string) error {
		if db.DB == nil {
			return nil
		}
		type GameRecordDB struct {
			GameID     string `gorm:"column:game_id;type:uuid"`
			TurnNumber int    `gorm:"column:turn_number"`
			GameState  string `gorm:"column:game_state;type:jsonb"`
		}
		record := GameRecordDB{
			GameID:     gameID,
			TurnNumber: turnNumber,
			GameState:  gameStateJSON,
		}
		return db.DB.Table("game_records").Create(&record).Error
	}

	gm.SaveFinalGame = func(gameID string, status string, currentTurn int, winnerID string, endedAt time.Time, players map[string]game.PlayerStateForSave) error {
		if db.DB == nil {
			return nil
		}

		var gameModel models.Game
		result := db.DB.Where("id = ?", gameID).First(&gameModel)
		if result.Error == nil {
			gameModel.Status = status
			gameModel.CurrentTurn = currentTurn
			gameModel.WinnerID = winnerID
			gameModel.EndedAt = endedAt
			db.DB.Save(&gameModel)
		}

		for pid, player := range players {
			var gp models.GamePlayer
			result := db.DB.Where("game_id = ? AND player_id = ?", gameID, pid).First(&gp)
			if result.Error == nil {
				gp.Score = player.Score
				gp.Eliminated = player.Eliminated
				db.DB.Save(&gp)
			}
		}

		return nil
	}
}
