FROM golang:alpine

LABEL version = "1.0"

LABEL maintainer="Dzafriel Zulkiflee <dzaf@vectolabs.com>"

WORKDIR $GOPATH/src/gitlab.com/verdant-sso

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 9992

EXPOSE 3000

CMD [ "cmd" ]

