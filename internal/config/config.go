package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Settings struct {
	Host                   string
	Port                   int
	RetrievalPeriodSeconds time.Duration
	SlackUrl               string
	SlackAuthToken         string
}

func NewSettings() (*Settings, error) {
	slackAuthToken := os.Getenv("SLACK_AUTH_TOKEN")
	if slackAuthToken == "" {
		return nil, errors.New("slack auth token is required, please set the SLACK_AUTH_TOKEN environment variable")
	}

	port := os.Getenv("SRES_PORT")
	if port == "" {
		port = "888"
	}

	nport, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("unable to parse SRES_PORT: %v", err)
	}

	retrievalPeriod := os.Getenv("RETRIEVAL_PERIOD_SECONDS")
	if retrievalPeriod == "" {
		retrievalPeriod = "60"
	}

	nRetrievalPeriod, err := strconv.Atoi(retrievalPeriod)
	if err != nil {
		return nil, fmt.Errorf("unable to parse RETRIEVAL_PERIOD_SECONDS: %v", err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	slackHost := os.Getenv("SLACK_HOST")
	if slackHost == "" {
		slackHost = "https://slack.com/api/emoji.list"
	}

	config := &Settings{host, nport, time.Duration(nRetrievalPeriod), slackHost, slackAuthToken}

	return config, nil
}
