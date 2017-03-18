FROM golang:1.6.2-wheezy

RUN mkdir -p /go/src/teamrunner/team

COPY clients/go/teamrunner.go /go/src/teamrunner/teamrunner.go
COPY clients/go/yielder /go/src/teamrunner/yielder
COPY clients/go/api /go/src/teamrunner/api
COPY clients/go/ids /go/src/teamrunner/ids
COPY clients/go/turninformer /go/src/teamrunner/turninformer
COPY clients/go/gameobjects /go/src/teamrunner/gameobjects
COPY clients/go/run.sh /go/src/teamrunner/run.sh
COPY vendor/git.apache.org /go/src/git.apache.org
COPY sourcemanager/sourcemanager /go/src/teamrunner/sourcemanager

WORKDIR /go/src/teamrunner
EXPOSE 9000

CMD ["./run.sh"]
