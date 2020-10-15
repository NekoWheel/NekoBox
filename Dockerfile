FROM alpine:latest

ADD . /home/app/
WORKDIR /home/app

RUN chmod 777 /home/app/NekoBox

ENTRYPOINT ["./NekoBox"]
EXPOSE 8080