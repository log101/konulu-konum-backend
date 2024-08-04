FROM golang:1.22.3-alpine AS builder

WORKDIR /usr/src/app

RUN apk update

RUN apk add build-base vips vips-dev

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=1 go build -v -o konulukonum

RUN go test -v ./...

FROM alpine
WORKDIR /root/app

RUN apk update

RUN apk add vips

COPY --from=builder /usr/src/app .

CMD ["./konulukonum"]
