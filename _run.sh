docker run -itd \
--privileged=true \
-v /etc/localtime:/etc/localtime:ro \
--name tgbot -p 5000:5000 tgbot:v1

