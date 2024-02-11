# swagger update
swagger-update-auth:
	@cd ./auth-app && swag init --parseDependency

swagger-update-chat:
	@cd ./chat-app && swag init --parseDependency

swagger-update-user:
	@cd ./user-app && swag init --parseDependency

swagger-update-game:
	@cd ./game-app && swag init --parseDependency
	
swagger-update-media:
	@cd ./media-app && swag init --parseDependency

swagger-update: swagger-update-auth swagger-update-chat swagger-update-user swagger-update-game

# migrating
migrate-auth:
	@echo "Migrating Auth app..."
	@cd ./auth-app && go run . migration up

migrate-chat:
	@echo "Migrating Chat app..."
	@cd ./chat-app && go run . migration up

migrate-game:
	@echo "Migrating Game app..."
	@cd ./game-app && go run . migration up

migrate-media:
	@echo "Migrating Media app..."
	@cd ./media-app && go run . migration up

migrate-user:
	@echo "Migrating User app..."
	@cd ./user-app && go run . migration up

migrate: migrate-auth migrate-chat migrate-game migrate-media migrate-user

# seeding database

seed-chat:
	@echo "Seeding Chat app..."
	@cd ./chat-app && go run . db seed

seed-user:
	@echo "Seeding User app..."
	@cd ./user-app && go run . db seed

seed: seed-chat seed-user



refactor: migrate seed

# Define variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PID_FILE=pids.txt

BINARY_AUTH_NAME=auth
BINARY_CHAT_NAME=chat
BINARY_GAME_NAME=game
BINARY_MEDIA_NAME=media
BINARY_USER_NAME=user
BINARY_GATEWAY_NAME=api-gateway

# Define targets
run: build run-apps

t:
	cd ./$(BINARY_GATEWAY_NAME)/bin && ./$(BINARY_GATEWAY_NAME) & echo $$! >> $(PID_FILE)

build:
	@echo "Building Golang applications..."
	@cd ./$(BINARY_AUTH_NAME)-app && $(GOBUILD) -o bin/$(BINARY_AUTH_NAME)  
	@cd ./$(BINARY_CHAT_NAME)-app && $(GOBUILD) -o bin/$(BINARY_CHAT_NAME)  
	@cd ./$(BINARY_GAME_NAME)-app && $(GOBUILD) -o bin/$(BINARY_GAME_NAME) 
	@cd ./$(BINARY_MEDIA_NAME)-app && $(GOBUILD) -o bin/$(BINARY_MEDIA_NAME)  
	@cd ./$(BINARY_USER_NAME)-app && $(GOBUILD) -o bin/$(BINARY_USER_NAME)   
	@cd ./$(BINARY_GATEWAY_NAME) && $(GOBUILD) -o bin/$(BINARY_GATEWAY_NAME) ./cmd

run-apps:
	@echo "Running Golang applications..."
	@cd ./$(BINARY_AUTH_NAME)-app && ./bin/$(BINARY_AUTH_NAME) serve & echo $$! >> $(PID_FILE)
	@cd ./$(BINARY_CHAT_NAME)-app && ./bin/$(BINARY_CHAT_NAME) serve & echo $$! >> $(PID_FILE)
	@cd ./$(BINARY_GAME_NAME)-app && ./bin/$(BINARY_GAME_NAME) serve & echo $$! >> $(PID_FILE)
	@cd ./$(BINARY_MEDIA_NAME)-app && ./bin/$(BINARY_MEDIA_NAME) serve & echo $$! >> $(PID_FILE)
	@cd ./$(BINARY_USER_NAME)-app && ./bin/$(BINARY_USER_NAME) serve & echo $$! >> $(PID_FILE)
	@cd ./$(BINARY_GATEWAY_NAME)/bin && ./$(BINARY_GATEWAY_NAME) & echo $$! >> $(PID_FILE)

stop:
	@echo "Stopping Golang applications..."
	@while read -r pid; do \
		kill "$$pid" || true; \
	done < $(PID_FILE)
	rm -f $(PID_FILE)