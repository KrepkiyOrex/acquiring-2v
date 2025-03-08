# ==================== for prod =============================
# FROM golang:1.24-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . ./

# WORKDIR /app/cmd

# RUN go build -v -o /main .

# # Этап запуска
# FROM alpine:latest

# # Копируем бинарник из этапа сборки
# COPY --from=builder /main /main

# # Даем права на выполнение
# RUN chmod +x /main

# # Запускаем приложение
# CMD ["/main"]

# ==================== for dev (CompileDaemon) =============================
FROM golang:1.24-alpine

RUN go install github.com/githubnemo/CompileDaemon@latest

# RUN apt-get update && apt-get install -y \
#     curl \
#     procps \
#     net-tools \
#     iproute2 \
#     dnsutils \
#     procps \
#     jq \
#     nano \
#     htop \
#     wget \
#     && rm -rf /var/lib/apt/lists/*

# WORKDIR /app/cmd
WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["CompileDaemon", "--build", "go build -o /app/cmd/main ./cmd", "--command", "/app/cmd/main"]
