-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
    id         VARCHAR(255) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255),
    avatar_url TEXT,
    github_id  BIGINT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ゲームデータテーブル
CREATE TABLE IF NOT EXISTS game_data (
    user_id                VARCHAR(255) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    main_character_id      VARCHAR(255) NOT NULL,
    play_time             INTEGER DEFAULT 0,
    stage                  INTEGER DEFAULT 1,
    last_commit_checked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス（検索を高速化）
CREATE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id);