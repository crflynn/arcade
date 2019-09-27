FROM golang:1.13.0

ENV GIN_MODE "release"

ENV ARCADE_PORT "6060"
ENV ARCADE_DOCROOT "/tmp/docs"
ENV ARCADE_USERNAME "admin"
ENV ARCADE_PASSWORD "password"

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o app

CMD ["./app"]