DOCKER_COMPOSE_FILE := docker-compose.yml
APP_NAME := redditclone
GO_CMD := go run cmd/$(APP_NAME)/main.go

up:
	@echo "Starting Docker services..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@echo "Stopping Docker services..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down

run:
	@echo "Running Go application..."
	$(GO_CMD)

mock:
	mockgen -source=./internal/repository/comment_repository.go -destination=./internal/mocks/mock_repo_comment.go -package=mocks
	mockgen -source=./internal/repository/post_repository.go -destination=./internal/mocks/mock_repo_post.go -package=mocks
	mockgen -source=./internal/repository/session_repository.go -destination=./internal/mocks/mock_repo_session.go -package=mocks
	mockgen -source=./internal/repository/user_repository.go -destination=./internal/mocks/mock_repo_user.go -package=mocks
	mockgen -source=./internal/service/service.go -destination=./internal/mocks/mock_service.go -package=mocks
	mockgen -source=./internal/token/token.go -destination=./internal/mocks/mock_token.go -package=mocks
	

.PHONY: up down run mock