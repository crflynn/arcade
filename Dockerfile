FROM golang:1.13.0

ENV GIN_MODE "release"

ENV GOLYGLOT_PORT "6060"
ENV GOLYGLOT_DOCROOT "/tmp/docs"
ENV GOLYGLOT_USERNAME "admin"
ENV GOLYGLOT_PASSWORD "password"

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o app

CMD ["./app"]