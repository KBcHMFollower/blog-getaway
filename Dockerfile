FROM golang:1.22.4-alpine

WORKDIR /app

COPY . .

RUN go build -o getaway cmd/init.go

CMD ["./getaway", "run-gw", "--cfg=config/prod.yaml"]