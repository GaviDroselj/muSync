FROM golang:1.25.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o musync cmd/main.go


FROM alpine

WORKDIR /app/

RUN apk add --no-cache ffmpeg

COPY --from=build /app/musync .

ENTRYPOINT [ "/app/musync" ]