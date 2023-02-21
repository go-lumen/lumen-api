FROM golang:latest

ENV SAM_ENV prod
ENV GIN_MODE release
ENV GO111MODULE on

RUN mkdir -p /var/www/uploads

WORKDIR /app
COPY go.mod .
COPY go.sum .é
RUN go mod download
COPY . .

RUN go version
RUN go build -o main

CMD ["/app/main"]