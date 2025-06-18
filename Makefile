#Переменные
APP_NAME = calc_restapi
DOCKER_IMAGE = calc_restapi
GO_FILES = $(shell find . -type f -name '*.go')
PORT = 8080

.PHONY: all build run test clean docker-build docker-run

all: build

#Сборка проекта
build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) ./cmd/server

#Запуск
run: build
	@echo "Starting..."
	@./bin/$(APP_NAME)

#Docker
docker-build:
	@echo "Docker building..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Starting docker..."
	@docker run -d -p $(PORT):$(PORT) --name $(APP_NAME) $(DOCKER_IMAGE)

docker-stop:
	@docker stop $(APP_NAME)

tidy:
	@go mod tidy -v

#Удаление созданного бинарника
clean:
	@echo "Cleaning..."
	@rm -rf bin/