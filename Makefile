swagger-update-auth:
	@cd ./auth-app && swag init --parseDependency

swagger-update-chat:
	@cd ./chat-app && swag init --parseDependency

swagger-update-user:
	@cd ./user-app && swag init --parseDependency

swagger-update-game:
	@cd ./game-app && swag init --parseDependency


swagger-update: swagger-update-auth swagger-update-chat swagger-update-user swagger-update-game