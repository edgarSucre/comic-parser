FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod main.go api ./

RUN go mod download && go mod verify
RUN go build -v ./...

FROM alpine:latest AS production

WORKDIR /
COPY --from=build /comic-parser /comic-parser
COPY public public

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT [ "/comic-parser" ]
