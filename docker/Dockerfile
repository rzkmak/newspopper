FROM golang:alpine as builder
RUN apk add --no-cache ca-certificates make tzdata
COPY . /app
RUN cd /app && \
    go get -d -v && \
    go mod download && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o newspopper .


FROM alpine:latest
LABEL maintainer="Anibal Aguila <aa@2bm.io>"

RUN mkdir /app
WORKDIR /app
VOLUME /data
COPY --from=builder /app/newspopper /app/
RUN apk add --no-cache ca-certificates tzdata tini
USER nobody
EXPOSE 8080
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/newspopper"]