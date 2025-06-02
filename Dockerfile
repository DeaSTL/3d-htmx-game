FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./game


FROM ubuntu:25.04 AS run-stage

WORKDIR /app
COPY --from=build-stage /app/maps ./maps
COPY --from=build-stage /app/public ./public
COPY --from=build-stage /app/game ./game

EXPOSE 3000


ENTRYPOINT ["sh","-c","./game"]
