FROM golang:1.24.5

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.3.924

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN templ generate

RUN go build -o ./game

ENTRYPOINT ["sh","-c","./game"]
