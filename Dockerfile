FROM golang:1.15.10-alpine
COPY . /go/src/go-spt-bot
WORKDIR /go/src/go-spt-bot
RUN go build -o /bin/spt-bot ./cmd/bot
ENV SPT_CFG_PATH=/data
ENTRYPOINT ["/bin/spt-bot"]
