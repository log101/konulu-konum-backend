FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o konulukonum

RUN go test -v ./...

FROM alpine
WORKDIR /root/app
COPY --from=builder /usr/src/app .

CMD ["./konulukonum"]
