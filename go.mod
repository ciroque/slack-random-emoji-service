module github.com/ciroque/slack-random-emoji-service

go 1.19

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	unit.nginx.org/go v0.0.0-20221219150355-8513083ac6d4 // indirect
)

replace (
	github.com/ciroque/slack-random-emoji-service/internal/config => ./internal/config
	github.com/ciroque/slack-random-emoji-service/internal/data => ./internal/data
	github.com/ciroque/slack-random-emoji-service/internal/data/sources => ./internal/data/sources
	github.com/ciroque/slack-random-emoji-service/internal/metrics => ./internal/metrics
	github.com/ciroque/slack-random-emoji-service/internal/server => ./internal/server
)
