FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/controlplane ./cmd/api

FROM debian:bookworm-slim

COPY --from=build /bin/controlplane /usr/local/bin/controlplane

EXPOSE 8080

CMD ["controlplane"]
