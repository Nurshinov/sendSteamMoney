FROM golang:alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /qiwi_web

FROM golang:alpine
COPY --from=build /app/qiwi_web /usr/bin/qiwi_web
CMD ["/usr/bin/qiwi_web"]