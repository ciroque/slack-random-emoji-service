package main

import (
	"fmt"
	"github.com/ciroque/slack-random-emoji-service/internal/config"
	"github.com/ciroque/slack-random-emoji-service/internal/data"
	"github.com/ciroque/slack-random-emoji-service/internal/data/sources"
	"github.com/ciroque/slack-random-emoji-service/internal/metrics"
	"github.com/ciroque/slack-random-emoji-service/internal/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopRetrieverChannel := make(chan bool)
	defer close(stopRetrieverChannel)

	abortChannel := make(chan string)
	defer close(abortChannel)

	emoUpdateChannel := make(chan *[]data.Emo)
	defer close(emoUpdateChannel)

	settings, error := config.NewSettings()
	if error != nil {
		logrus.Fatalf("Error creating configuration settings: %v", error)
	}

	_ = fmt.Sprintf("[[[[[ %v", settings)

	metricClient := metrics.NewMetrics()

	var emos *[]data.Emo

	httpServer := server.Server{
		AbortChannel:     abortChannel,
		Logger:           logrus.NewEntry(logrus.New()),
		Emos:             emos,
		EmoUpdateChannel: emoUpdateChannel,
		Settings:         settings,
		Metrics:          &metricClient,
	}

	slackEmoRetriever := sources.SlackRetriever{
		EmoUpdateChannel: emoUpdateChannel,
		Settings:         settings,
		StopChannel:      stopRetrieverChannel,
		Metrics:          &metricClient,
	}

	go httpServer.Run()
	go httpServer.HandleUpdates()
	go slackEmoRetriever.Run()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)

	select {
	case <-sigTerm:
		{
			stopRetrieverChannel <- true
			logrus.Info("Exiting per SIGTERM")
		}
	case err := <-abortChannel:
		{
			stopRetrieverChannel <- true
			logrus.Error(err)
		}
	}
}
