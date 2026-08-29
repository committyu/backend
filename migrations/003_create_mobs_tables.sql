-- モブテーブル
CREATE TABLE IF NOT EXISTS mobs (
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    hp          INTEGER DEFAULT 0,
    atk         INTEGER DEFAULT 0,
    matk        INTEGER DEFAULT 0,
    def         INTEGER DEFAULT 0,
    mdef        INTEGER DEFAULT 0,
    agi         INTEGER DEFAULT 0,
    xp          INTEGER DEFAULT 0,
);