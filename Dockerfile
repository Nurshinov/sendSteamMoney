FROM golang:alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -tags musl -o /tg_bot
CMD ["/tg_bot"]