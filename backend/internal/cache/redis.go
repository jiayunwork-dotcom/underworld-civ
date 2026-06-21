package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"underworld-civ/internal/config"
	"underworld-civ/internal/game"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func Init(cfg *config.Config) error {
	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	_, err := RDB.Ping(ctx).Result()
	return err
}

func GetGameState(gameID string) (*game.GameState, error) {
	key := fmt.Sprintf("game:%s", gameID)
	data, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var state game.GameState
	err = json.Unmarshal([]byte(data), &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func SetGameState(state *game.GameState) error {
	key := fmt.Sprintf("game:%s", state.ID)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return RDB.Set(ctx, key, data, 0).Err()
}

func DeleteGameState(gameID string) error {
	key := fmt.Sprintf("game:%s", gameID)
	return RDB.Del(ctx, key).Err()
}

func GetRoom(gameID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("room:%s", gameID)
	data, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}

	return result, nil
}

func SetRoomField(gameID, field, value string) error {
	key := fmt.Sprintf("room:%s", gameID)
	return RDB.HSet(ctx, key, field, value).Err()
}

func DeleteRoom(gameID string) error {
	key := fmt.Sprintf("room:%s", gameID)
	return RDB.Del(ctx, key).Err()
}

func AddPlayerToRoom(gameID, playerID, username string) error {
	key := fmt.Sprintf("room:%s:players", gameID)
	return RDB.HSet(ctx, key, playerID, username).Err()
}

func RemovePlayerFromRoom(gameID, playerID string) error {
	key := fmt.Sprintf("room:%s:players", gameID)
	return RDB.HDel(ctx, key, playerID).Err()
}

func GetRoomPlayers(gameID string) (map[string]string, error) {
	key := fmt.Sprintf("room:%s:players", gameID)
	return RDB.HGetAll(ctx, key).Result()
}

type PlayerTechPublic struct {
	Techs      []string `json:"techs"`
	Current    string   `json:"current_research"`
	PlayerID   string   `json:"player_id"`
	Username   string   `json:"username"`
}

func SyncPlayerTechs(gameID string, playerID string, techs []string, currentResearch string, username string) error {
	key := fmt.Sprintf("game:%s:techs", gameID)
	data := PlayerTechPublic{
		Techs:    techs,
		Current:  currentResearch,
		PlayerID: playerID,
		Username: username,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return RDB.HSet(ctx, key, playerID, jsonData).Err()
}

func GetAllPlayerTechs(gameID string) (map[string]PlayerTechPublic, error) {
	key := fmt.Sprintf("game:%s:techs", gameID)
	rawData, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]PlayerTechPublic)
	for pid, raw := range rawData {
		var ptp PlayerTechPublic
		if err := json.Unmarshal([]byte(raw), &ptp); err == nil {
			result[pid] = ptp
		}
	}
	return result, nil
}

func DeletePlayerTechs(gameID string) error {
	key := fmt.Sprintf("game:%s:techs", gameID)
	return RDB.Del(ctx, key).Err()
}
