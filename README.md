# 起動コマンド
## 初回・ビルド時
Bash
docker compose up --build

## 2回目以降（バックグラウンド起動）
Bash
docker compose up -d


## 開発用コマンド
## ライブラリの追加
Bash
docker compose exec api go get <ライブラリ名>
docker compose exec api go mod tidy
## データベース(PostgreSQL)への接続確認
Bash
docker compose exec db psql -U user -d game_db
## Redisの動作確認
Bash
docker compose exec redis redis-cli ping
"PONG" と返ってくれば正常です

# 技術スタック
Language: Go 1.25.1

Framework: Echo

ORM: GORM

Hot Reload: Air

DB: PostgreSQL / Redis

Worker: Asynq