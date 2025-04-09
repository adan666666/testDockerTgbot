#docker stop tgbot
docker rm -f tgbot
docker rmi -f tgbot
docker build -t tgbot:v1 .