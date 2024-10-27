FROM golang:1.23.2 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY ./vendor/ ./vendor/
# RUN go mod download -x \
#  && go mod verify

COPY main.go ./
COPY internal/ ./internal/

RUN GOMAXPROCS=4 CGO_ENABLED=0 \
    go build -a -mod=vendor -o my-controller main.go

FROM busybox
WORKDIR /app
COPY --from=builder /usr/src/app/my-controller ./
CMD ["./my-controller"]
