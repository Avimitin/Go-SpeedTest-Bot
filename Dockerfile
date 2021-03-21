FROM golang:latest
COPY . /go/src/go-spt-bot
WORKDIR /go/src/go-spt-bot
RUN go build -o /bin/spt-bot ./cmd/bot
ENTRYPOINT ["/bin/spt-bot"]
