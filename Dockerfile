FROM golang:1.23

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /bike-shop

EXPOSE 8070

CMD ["/bike-shop"]
