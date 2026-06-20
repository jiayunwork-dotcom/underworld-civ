CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    max_players INTEGER NOT NULL DEFAULT 6,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    current_turn INTEGER DEFAULT 0,
    max_turns INTEGER DEFAULT 50,
    winner_id UUID REFERENCES players(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS game_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    player_id UUID REFERENCES players(id),
    race VARCHAR(30) NOT NULL,
    color VARCHAR(20) NOT NULL,
    is_host BOOLEAN DEFAULT false,
    score INTEGER DEFAULT 0,
    eliminated BOOLEAN DEFAULT false,
    UNIQUE(game_id, player_id)
);

CREATE TABLE IF NOT EXISTS game_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    turn_number INTEGER NOT NULL,
    game_state JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_games_status ON games(status);
CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id);
CREATE INDEX IF NOT EXISTS idx_game_records_game_id ON game_records(game_id);
