FROM golang:1.6.2-wheezy

RUN mkdir -p /go/src/botrunner/bot

COPY clients/go/botrunner.go /go/src/botrunner/botrunner.go
COPY clients/go/yielder /go/src/botrunner/yielder
COPY clients/go/api /go/src/botrunner/api
COPY clients/go/ids /go/src/botrunner/ids
COPY clients/go/turninformer /go/src/botrunner/turninformer
COPY clients/go/gameobjects /go/src/botrunner/gameobjects
COPY clients/go/runbot.sh /go/src/botrunner/runbot.sh
COPY vendor/git.apache.org /go/src/git.apache.org
COPY sourcemanager/sourcemanager /go/src/botrunner/sourcemanager

WORKDIR /go/src/botrunner
EXPOSE 9000

CMD ["./runbot.sh"]
