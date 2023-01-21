FROM golang:1.19-alpine

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN go build -o ./build/service ./cmd/main/.

CMD [ "./build/service" ]
