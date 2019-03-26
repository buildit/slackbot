############################
# STAGE 1 build Golang executable binary
############################

FROM golang:1.12.1-alpine AS builder
#https://github.com/golang/go/issues/28065
ENV CGO_ENABLED=0
RUN apk update && apk add --no-cache ca-certificates openssl openssh curl bash git && update-ca-certificates
WORKDIR $GOPATH/src/github.com/buildit/slackbot
ADD . $GOPATH/src/github.com/buildit/slackbot/
RUN go get -d -v ./...
RUN go get -u github.com/jstemmer/go-junit-report
RUN mkdir /go/TestResults
RUN go test -v ./... | go-junit-report > /go/TestResults/TestReport.xml
RUN GOOS=linux GOARCH=amd64 go build -v -o /go/bin/bot-server ./cmd/bot-server/main.go

############################
# STAGE 2 Upload test results to Azure
############################

FROM microsoft/dotnet:latest AS tester

ARG STORAGE_ACCT_URL
ARG STORAGE_ACCT_KEY
ARG BUILD_NUMBER

RUN apt-get update && apt-get -y install rsync --no-install-recommends apt-utils
RUN mkdir /tmp/azcopy && \
    wget -O /tmp/azcopy/azcopy.tar.gz https://aka.ms/downloadazcopyprlinux &&  \
    tar -xf /tmp/azcopy/azcopy.tar.gz -C /tmp/azcopy &&  \
    /tmp/azcopy/install.sh

RUN rm -rf /tmp/azcopy

COPY --from=builder /go/TestResults/TestReport.xml ./TestReport.xml

RUN azcopy \
     --source ./TestReport.xml \
     --destination "${STORAGE_ACCT_URL}/TestReport_${BUILD_NUMBER}.xml" \
     --dest-key "${STORAGE_ACCT_KEY}"


############################
# STAGE 3 Build small image with only binary
############################

FROM alpine:3.9 AS final
EXPOSE 4390
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=builder /go/bin/bot-server /go/bin/bot-server
ENTRYPOINT ["/go/bin/bot-server"]