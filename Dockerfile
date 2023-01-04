############################
# STEP 1 build executable binary
############################
FROM nginx/unit:1.29.0-go1.19 AS builder

# Install necessary build packages
RUN apt-get update && apt-get install -y git build-essential libpcre3-dev

WORKDIR ~
RUN git clone https://github.com/nginx/unit.git && \
    cd unit && \
    ./configure \
      --prefix=/usr/local \
      --bindir=/usr/local/bin \
      --sbindir=/usr/local/sbin \
      --modules=/usr/local/lib/unit/modules \
      --state=/usr/local/var/unit/state \
      --pid=/usr/local/var/run/unit.pid \
      --log=/usr/local/var/log/unit.log \
      --user="$(id -un)" \
      --group="$(id -gn)" \
      --control=unix:/usr/local/var/run/control.unit.sock && \
    ./configure go && \
    make && \
    make install && \
    make go-install

WORKDIR $GOPATH/src/slack-random-emoji-service
COPY . .

# Fetch dependencies.
# Using go install.
RUN go install ./cmd/slack-random-emoji-service/main.go

# Build the binary.
RUN go build -o /tmp/slack-random-emoji-service ./cmd/slack-random-emoji-service/main.go

# Stage the configuration file
COPY ./config/unit.json /tmp/unit.json

############################
# STEP 2 build a small image
############################
FROM nginx/unit:1.29.0-go1.19

WORKDIR /opt/slack-random-emoji-service/bin

# Copy our static executable.
COPY --from=builder /tmp/slack-random-emoji-service ./slack-random-emoji-service

# Copy our configuation.
COPY --from=builder /tmp/unit.json /docker-entrypoint.d/unit.json

# Run the hello binary.

EXPOSE 20120

# docker run -d -p 20120:20120 slack-random-emoji-service-unit