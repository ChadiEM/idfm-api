FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o idfm -ldflags "-s -w" /app/cmd/idfm

FROM gcr.io/distroless/static-debian12

COPY --from=build /app/idfm /

USER nonroot

ENTRYPOINT ["/idfm"]