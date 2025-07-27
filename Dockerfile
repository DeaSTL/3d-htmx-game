FROM golang:1.24.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@v0.3.924

RUN templ generate

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./game


FROM ubuntu:25.04 AS run-stage

WORKDIR /app
COPY --from=build-stage /app/maps ./maps
COPY --from=build-stage /app/public ./public
COPY --from=build-stage /app/game ./game
COPY --from=build-stage /app/prefabs ./prefabs

EXPOSE 3000


ENTRYPOINT ["sh","-c","./game"]
