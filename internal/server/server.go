package server

import (
	"encoding/json"
	"fmt"
	"github.com/ciroque/slack-random-emoji-service/internal/config"
	"github.com/ciroque/slack-random-emoji-service/internal/data"
	"github.com/ciroque/slack-random-emoji-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	unit "unit.nginx.org/go"
)

type Server struct {
	AbortChannel     chan<- string
	Logger           *logrus.Entry
	Emos             *[]data.Emo
	EmoUpdateChannel <-chan *[]data.Emo
	Settings         *config.Settings
	Metrics          *metrics.Metrics
}

func (server *Server) Run() {
	randomEmojiRequestHandler := promhttp.InstrumentHandlerCounter(
		server.Metrics.RandomEmoRequestCount,
		promhttp.InstrumentHandlerDuration(server.Metrics.RandomEmoRequestDurations,
			http.HandlerFunc(server.ServeRandomEmoji)))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", randomEmojiRequestHandler)

	address := fmt.Sprintf("%s:%d", server.Settings.Host, server.Settings.Port)
	server.Logger.Info("Listening on ", address)
	err := unit.ListenAndServe(address, nil)
	if err != nil {
		server.AbortChannel <- err.Error()
	}
}

func (server *Server) ServeRandomEmoji(writer http.ResponseWriter, _ *http.Request) {
	length := len(*server.Emos)
	index := rand.Intn(length)

	response := SlackEmojiResponse{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf(":%s:", (*server.Emos)[index].Name),
		Attachments:  []map[string]string{},
	}

	bytes, err := json.Marshal(&response)
	if err != nil {
		server.Logger.Warnf("Error responding to request %#v", err)
	}

	writer.Header().Add("Content-Type", "application/json")
	_, err = fmt.Fprintf(writer, "%s", bytes)
	if err != nil {
		server.Logger.Warnf("Error responding to request %#v", err)
	}
}

func (server *Server) HandleUpdates() {
	for updatedEmos := range server.EmoUpdateChannel {
		server.Emos = updatedEmos
	}
}
