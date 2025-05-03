FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg

RUN go build -o idfm /app/cmd/idfm

FROM alpine:3.21.3

WORKDIR /

COPY --from=build /app/idfm .

USER nobody

ENTRYPOINT [ "/idfm" ]
