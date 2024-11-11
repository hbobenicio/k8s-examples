FROM golang:1.23.2 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY ./vendor/ ./vendor/
# RUN go mod download -x \
#  && go mod verify

COPY main.go ./
COPY internal/ ./internal/

RUN GOMAXPROCS=4 CGO_ENABLED=0 \
    go build -a -ldflags="-s -w" -mod=vendor -o my-controller main.go

FROM scratch
WORKDIR /app
ENV GOMAXPROCS=4
COPY --from=builder /usr/src/app/my-controller ./
CMD ["./my-controller"]
