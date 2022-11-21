# builder docker image
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# build runtime image

FROM alpine:latest

RUN mkdir /app

COPY --chown=brokerUser:brokerUser --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]
