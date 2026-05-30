# Postman 用 API 定義書

## 基本情報

- Base URL: `http://localhost:8080`
- API Prefix: `/api`
- Content-Type: `application/json`

Postman の Environment に以下を設定してください。

| 変数名 | 例 | 説明 |
| --- | --- | --- |
| `baseUrl` | `http://localhost:8080` | バックエンドの URL |
| `githubCode` | `xxxxxxxx` | GitHub OAuth で発行される一時 code |
| `jwtToken` | `eyJ...` | ログイン API が返す JWT |

## 認証の確認手順

1. ブラウザで GitHub 認可 URL を開きます。

```text
https://github.com/login/oauth/authorize?client_id=<GITHUB_CLIENT_ID>&redirect_uri=<GITHUB_REDIRECT_URL>&scope=user:email
```

2. GitHub 認証後、`GITHUB_REDIRECT_URL` にリダイレクトされます。
3. リダイレクト先 URL の `code` クエリパラメータをコピーします。
4. コピーした値を Postman の `githubCode` に設定します。
5. `GET /api/auth/login?code={{githubCode}}` を実行します。
6. レスポンスの `token` を Postman の `jwtToken` に設定します。
7. 認証が必要な API では `Authorization: Bearer {{jwtToken}}` を付けて実行します。

補足: `GITHUB_REDIRECT_URL` を `http://localhost:8080/api/auth/callback` にしている場合は、ブラウザで GitHub 認証後にそのままローカル API が叩かれ、JSON で JWT が返ります。この場合は表示された `token` を Postman の `jwtToken` に設定すれば確認できます。

## API 一覧

### GitHub ログイン

GitHub OAuth の一時 code を使って GitHub access token を取得します。取得した access token は暗号化して DB に保存し、このアプリ用の JWT を返します。

- Method: `GET`
- URL: `{{baseUrl}}/api/auth/login?code={{githubCode}}`
- Auth: なし

Query Parameters:

| Name | Required | 説明 |
| --- | --- | --- |
| `code` | Yes | GitHub OAuth の一時 code |

成功レスポンス: `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

エラーレスポンス例:

```json
{
  "error": "code is required"
}
```

```json
{
  "error": "internal server error"
}
```

Postman の Tests に以下を設定すると、レスポンスの `token` を `jwtToken` に自動保存できます。

```javascript
const body = pm.response.json();
if (body.token) {
  pm.environment.set("jwtToken", body.token);
}
```

### GitHub コールバック

`/api/auth/login` と同じ処理です。`GITHUB_REDIRECT_URL` をバックエンドの callback API に向ける場合に使います。

- Method: `GET`
- URL: `{{baseUrl}}/api/auth/callback?code={{githubCode}}`
- Auth: なし

成功レスポンス: `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 認証が必要な API

以下の API は Bearer Token が必要です。

Postman の Authorization タブで以下を設定してください。

- Type: `Bearer Token`
- Token: `{{jwtToken}}`

または Headers に直接設定します。

| Name | Value |
| --- | --- |
| `Authorization` | `Bearer {{jwtToken}}` |

### 自分のユーザー情報取得

ログイン中のユーザー情報を取得します。

- Method: `GET`
- URL: `{{baseUrl}}/api/user/me`
- Auth: Bearer Token

成功レスポンス: `200 OK`

```json
{
  "id": "b388d3d3-04a5-465b-bb7f-f22707c03dd4",
  "github_name": "hsmt-T",
  "email": "",
  "avatar_url": "https://avatars.githubusercontent.com/u/208746076?v=4",
  "created_at": "2026-05-04T06:24:32Z"
}
```

エラーレスポンス例:

```json
{
  "error": "invalid token"
}
```

```json
{
  "error": "user not found"
}
```

### GitHub Commit 同期

DB に保存されている GitHub access token を使って GitHub API から PushEvent を取得し、ゲーム側の commit 数を更新します。

- Method: `GET`
- URL: `{{baseUrl}}/api/game/syncCommit`
- Auth: Bearer Token

成功レスポンス: `200 OK`

```json
{
  "message": "sync completed",
  "github_name": "hsmt-T",
  "checked_since": "2026-05-04T06:24:32Z",
  "checked_until": "2026-05-30T12:00:00Z",
  "matched_push_event_count": 2,
  "new_commit_count": 5,
  "total_commits": 15,
  "updated": true
}
```

エラーレスポンス例:

```json
{
  "error": "invalid token"
}
```

```json
{
  "error": "github access token not found"
}
```

```json
{
  "error": "github events api error: status 401"
}
```

## ローカル確認用メモ

`.env` に必要な値:

```env
DB_HOST=db
DB_USER=user
DB_PASSWORD=password
DB_NAME=game_db
DB_PORT=5432

GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URL=
GITHUB_TOKEN_ENCRYPTION_KEY=

JWT_SECRET=
```

`GITHUB_TOKEN_ENCRYPTION_KEY` は GitHub access token を暗号化するための秘密文字列です。任意の文字列を設定できますが、本番では推測されにくい長い値にしてください。

サーバー起動:

```bash
docker compose up --build
```

## Postman での確認順

1. `GET {{baseUrl}}/api/auth/login?code={{githubCode}}`
2. `GET {{baseUrl}}/api/user/me`
3. `GET {{baseUrl}}/api/game/syncCommit`
