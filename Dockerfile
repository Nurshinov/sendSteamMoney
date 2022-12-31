FROM golang:alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/qiwi_web

FROM golang:alpine
EXPOSE 8080
COPY --from=build /app/qiwi_web /usr/bin/qiwi_web
CMD ["/usr/bin/qiwi_web"]