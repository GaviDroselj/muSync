FROM golang:1.25.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o musync cmd/main.go

FROM golang:1.25.1 AS dependencies

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY scripts scripts
COPY dependencies dependencies

RUN go run scripts/download-dependencies.go

FROM denoland/deno

WORKDIR /app/

COPY --from=dependencies /root/.cache /root/.cache

COPY --from=build /app/musync .

CMD [ "/app/musync" ]