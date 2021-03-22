# BUILD
FROM golang:1.15.10-alpine AS build
COPY . /go/src/go-spt-bot
WORKDIR /go/src/go-spt-bot
RUN go build -o /bin/spt-bot ./cmd/bot

# RUN
FROM alpine:3
COPY --from=build /bin/spt-bot /bin/spt-bot
ENV SPT_CFG_PATH=/data
ENTRYPOINT ["/bin/spt-bot"]
