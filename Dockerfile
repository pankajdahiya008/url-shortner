FROM golang:1.23

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortner

EXPOSE 8080

CMD ["/url-shortner"]