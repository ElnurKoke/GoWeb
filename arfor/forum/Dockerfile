FROM golang:alpine AS builder


WORKDIR /app
COPY . .

# Скачивает гит без кэша && инициализирует гоу мод файл && создает бинарный файл форум  из cmd/main.go
RUN apk add --no-cache build-base && \
    CGO_ENABLED=1 go mod download && \
    CGO_ENABLED=1 go build -o forum cmd/main.go
# Берет второй раз . но на этот раз последнюю версию дистрибутива линукса
FROM alpine:latest

# АВТОР && ПРОЕКТ && ВРЕМЯ
LABEL authors="Aryn&Ayan" \
      project="Forum" \
      date="2412412"

# Создает новую папку app для FROM alpine:latest
WORKDIR /app

# Копирует бинарный файл
COPY --from=builder /app .

# # Удаляет гит
# RUN apk del git
EXPOSE 1111
# бинарный файл
CMD ["./forum"]
