FROM alpine:latest

RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

ADD NekoBox /home/app/NekoBox
WORKDIR /home/app

RUN chmod 777 /home/app/NekoBox

ENTRYPOINT ["./NekoBox"]
EXPOSE 8080
