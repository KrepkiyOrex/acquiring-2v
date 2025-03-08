.SILENT:

hello:
	echo "hello"

start:
	docker compose up --build


# stop docker compose
stop:
	docker compose down


# exec consumer_container
con-bash:
	docker exec -it consumer_cont /bin/sh


# rebuild consumer_container
con-bd:
	docker compose up --build consumer-app

# rebuild acquiring_container
app-bd:
	docker compose up --build app


# delete folder
del:
	sudo rm -rf ./internal
	sudo rm -rf ./consumer


# docker stop acquiring_container
# docker stop acquiring_container
# # docker stop consumer_container
# alias stc="docker stop consumer_container_container"