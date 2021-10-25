FROM golang:1.15-alpine3.14

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build

EXPOSE 8080

CMD ["./discrepancy"]
