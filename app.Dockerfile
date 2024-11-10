# Используем официальный Go образ
FROM golang:1.23.2-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Загружаем зависимости
RUN go mod tidy

ENV CONFIG_PATH="./config/config.yaml"

# Строим приложение
RUN go build -o main ./cmd/app/main.go

# Указываем команду для запуска
CMD ["./main"]