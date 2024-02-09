# swagger update
swagger-update-auth:
	@cd ./auth-app && swag init --parseDependency

swagger-update-chat:
	@cd ./chat-app && swag init --parseDependency

swagger-update-user:
	@cd ./user-app && swag init --parseDependency

swagger-update-game:
	@cd ./game-app && swag init --parseDependency

swagger-update: swagger-update-auth swagger-update-chat swagger-update-user swagger-update-game

# deploy
deploy-auth:
	liara deploy --app chess-auth --port 8001 --platform docker --dockerfile ./auth-app/deployments/docker/Dockerfile 

deploy-user:
	liara deploy --app chess-user --port 8004 --platform docker --dockerfile ./user-app/deployments/docker/Dockerfile