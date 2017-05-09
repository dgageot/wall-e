FROM golang:1.8.1-alpine as builder

COPY pkg .
RUN go build -a -tags "static_build netgo" -o app

FROM alpine:3.5

RUN apk add --update ca-certificates

ENV JENKINS_SERVER ""
ENV JENKINS_USER ""
EXPOSE 8080 8888

ENTRYPOINT ["/wall-e"]
COPY --from=builder /go/app /wall-e
