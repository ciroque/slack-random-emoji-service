############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git build-base
WORKDIR $GOPATH/src/slack-random-emoji-service
COPY . .

# Fetch dependencies.
# Using go get.
RUN go install ./cmd/slack-random-emoji-service/main.go

# Build the binary.
RUN go build -o /tmp/slack-random-emoji-service ./cmd/slack-random-emoji-service/main.go

############################
# STEP 2 build a small image
############################
FROM alpine:3.16.3

WORKDIR /opt/slack-random-emoji-service/bin

# Copy our static executable.
COPY --from=builder /tmp/slack-random-emoji-service ./slack-random-emoji-service

# Run the hello binary.
ENTRYPOINT ["/opt/slack-random-emoji-service/bin/slack-random-emoji-service"]
