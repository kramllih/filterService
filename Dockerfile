#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./cmd/

#final stage
FROM alpine:latest
ENV GIN_MODE=release
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/app/cmd/config.yml /config.yml
ENTRYPOINT /app
LABEL Name=filter Version=0.0.1
EXPOSE 8080
