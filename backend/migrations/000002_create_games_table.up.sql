CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    difficulty VARCHAR(20) NOT NULL,
    result VARCHAR(10),
    board_state JSONB,
    current_turn VARCHAR(1),
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE INDEX idx_games_user_id ON games(user_id);
CREATE INDEX idx_games_created_at ON games(created_at DESC);
