FROM golang:1.25.1

WORKDIR /app

# air install
RUN go install github.com/air-verse/air@latest

# go.modコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコピー
COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]