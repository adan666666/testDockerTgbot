#创建镜像
#docker stop tgbot
docker rm -f tgbot #删除容器
docker rmi -f tgbot:v1.0 #删除镜像
docker build -t tgbot:v1.0 . #构建镜像tgbot:v1.0

#创建并启动容器
docker run -itd \
--privileged=true \
-v /etc/localtime:/etc/localtime:ro \
--name tgbot -p 5000:5000 tgbot:v1.0

