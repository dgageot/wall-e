FROM golang:1.8.1-alpine as builder

COPY pkg .
RUN go build -a -tags "static_build netgo" -o app

FROM alpine:3.5

RUN apk add --update ca-certificates

ENV JENKINS_SERVER ""
ENV JENKINS_USER ""

ENTRYPOINT ["/wall-e"]
CMD ["8080"]
COPY --from=builder /go/app /wall-e
