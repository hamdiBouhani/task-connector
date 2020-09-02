FROM golang:1.12.3-alpine3.9 as builder

RUN apk add --no-cache --update git alpine-sdk


COPY . /go/src/gitlab.com/target-smart-data-ai-searsh/task-connector-be 
RUN go build -o /go/bin/extractor gitlab.com/target-smart-data-ai-searsh/task-connector-be/cmd

FROM alpine

ENV "GOPATH" "/go"

RUN apk add -u --no-cache tzdata ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/gitlab.com/target-smart-data-ai-searsh/task-connector-be/service-account-file.json /usr/local/bin/

ENV GOOGLE_APPLICATION_CREDENTIALS "/usr/local/bin/service-account-file.json"

COPY --from=builder /go/bin/extractor /usr/local/bin/

WORKDIR /
