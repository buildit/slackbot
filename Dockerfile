############################
# STAGE 1 build Golang executable binary
############################

FROM golang:1.10-alpine AS builder
ENV bin_dir=/go/bin/
RUN apk update && apk add --no-cache openssl openssh curl bash git openssh libcurl
WORKDIR $GOPATH/src/github.com/buildit/slackbot
ADD . $GOPATH/src/github.com/buildit/slackbot/
RUN go get -d -v ./...
RUN go build -v -o ${bin_dir}/bot-server.sh ./cmd/bot-server/main.go
RUN go test ./...

############################
# STAGE 2 Build small image with only binary
############################

FROM golang:1.10-alpine
ENV bin_dir=/go/bin/
COPY --from=builder ${bin_dir}/bot-server.sh /go/bin/bot-server.sh
RUN cd ${bin_dir} && ls
ENTRYPOINT ["/go/bin/bot-server.sh"]