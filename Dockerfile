FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod go.mod
RUN go mod download && go mod verify

COPY api api
COPY main.go main.go

RUN go build

# FROM alpine:latest AS production

# WORKDIR /app
# COPY --from=build /app/comic-parser .
COPY public public

EXPOSE 8000
CMD ["./comic-parser"]
