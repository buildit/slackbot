############################
# STAGE 1 build Golang executable binary
############################

FROM golang:1.12.5-alpine AS builder
#https://github.com/golang/go/issues/28065
ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN apk update && apk add --no-cache ca-certificates openssl openssh curl bash git && update-ca-certificates
RUN mkdir /slackbot && mkdir /slackbot/TestResults
WORKDIR /slackbot
COPY go.mod .
COPY go.sum .
RUN go mod download build
COPY . .
RUN go get -u github.com/jstemmer/go-junit-report
RUN go test -v ./... -tags unit_tests | go-junit-report > /slackbot/TestResults/TestReport.xml
RUN GOOS=linux GOARCH=amd64 go build -v -o /slackbot/bot-server ./main.go

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

COPY --from=builder /slackbot/TestResults/TestReport.xml ./TestReport.xml

RUN azcopy \
     --source ./TestReport.xml \
     --destination "${STORAGE_ACCT_URL}/TestReport_${BUILD_NUMBER}.xml" \
     --dest-key "${STORAGE_ACCT_KEY}" \
     --dest-type blob


############################
# STAGE 3 Build small image with only binary
############################

FROM alpine:3.9 AS final
EXPOSE 4390
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=builder /slackbot/bot-server /slackbot/bot-server
ENTRYPOINT ["/slackbot/bot-server"]
