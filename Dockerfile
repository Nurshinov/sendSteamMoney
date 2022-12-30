FROM golang:alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /qiwi_web
CMD ["/qiwi_web"]