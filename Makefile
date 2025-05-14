PROJECT_NAME := rtb

up:
	sudo docker compose -p $(PROJECT_NAME) up --force-recreate --no-deps --build -d

down:
	sudo docker compose -p $(PROJECT_NAME) down