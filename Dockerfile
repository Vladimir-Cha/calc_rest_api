# Используем официальный образ Go на основе Alpine
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Компилируем приложение
RUN go build -o /Calc_restapi ./cmd/server/

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /

# Копируем скомпилированное приложение из предыдущего образа
COPY --from=builder /Calc_restapi /Calc_restapi

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD [".//Calc_restapi"]