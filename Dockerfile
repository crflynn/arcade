# start with a builder image and move to scratch later
FROM golang:1.13.0-alpine as builder

# add the ca certificates
RUN apk update && \
    apk add --update --no-cache ca-certificates && \
    update-ca-certificates

# set the working directory
WORKDIR /build

# add dependency configs
ADD go.mod .
ADD go.sum .

# install dependencies
RUN go mod download && \
    go mod verify

# add the application code
ADD . .

# build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .


# use scratch for a minimal image
FROM scratch

ENV ARCADE_PORT "6060"
ENV ARCADE_DOCROOT "/docs"

# copy the certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

# copy our executable
COPY --from=builder /build/main /app/main

# set the application to release mode
ENV GIN_MODE="release"

# set the entrypoint to the executable
ENTRYPOINT ["/app/main"]
