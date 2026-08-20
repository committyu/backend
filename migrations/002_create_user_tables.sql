-- キャラクターテーブル
CREATE TABLE IF NOT EXISTS characters (
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    job         TEXT,
    hp          INTEGER DEFAULT 0,
    atk         INTEGER DEFAULT 0,
    matk        INTEGER DEFAULT 0,
    def         INTEGER DEFAULT 0,
    mdef        INTEGER DEFAULT 0,
    agi         INTEGER DEFAULT 0,
    luk         INTEGER DEFAULT 0,
    xp          INTEGER DEFAULT 0,
    user_id     VARCHAR(255) PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS id ON users(id);