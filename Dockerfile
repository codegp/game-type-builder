FROM gcr.io/codegamingplatform-146701/builderbase

RUN mkdir -p /go/src/github.com/codegp/game-type-builder
COPY vendor/ /go/src/github.com/codegp/game-type-builder/vendor

COPY thriftgenerator/ /go/src/github.com/codegp/game-type-builder/thriftgenerator
COPY sourcemanager/ /go/src/github.com/codegp/game-type-builder/sourcemanager
COPY docsreporter/ /go/src/github.com/codegp/game-type-builder/docsreporter
COPY buildgametype.sh /go/src/github.com/codegp/game-type-builder/buildgametype.sh
COPY clients/ /go/src/github.com/codegp/game-type-builder/clients
COPY game-runner/ /go/src/github.com/codegp/game-runner

WORKDIR /go/src/github.com/codegp/game-type-builder

CMD ["./buildgametype.sh"]
