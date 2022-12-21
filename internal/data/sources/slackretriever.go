package sources

import (
	"encoding/json"
	"fmt"
	"github.com/ciroque/slack-random-emoji-service/internal/config"
	"github.com/ciroque/slack-random-emoji-service/internal/data"
	"github.com/ciroque/slack-random-emoji-service/internal/metrics"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type SlackRetriever struct {
	EmoUpdateChannel chan<- *[]data.Emo
	Settings         *config.Settings
	StopChannel      <-chan bool
	Metrics          *metrics.Metrics
}

func (retriever *SlackRetriever) Run() {
	//fmt.Println(">>> ", retriever, "<<<", retriever.Settings, ">>>")
	periodic := time.NewTicker(time.Second * retriever.Settings.RetrievalPeriodSeconds)

	retriever.worker()

	go func() {
		for {
			select {
			case <-retriever.StopChannel:
				{
					logrus.Info("Shutting down Slack Retriever")
					return
				}
			case t := <-periodic.C:
				{
					logrus.Info("Tick at ", t)
					retriever.worker()
				}
			}
		}
	}()
}

func (retriever *SlackRetriever) worker() {
	retriever.Metrics.EmoRetrievalCount.Inc()
	started := time.Now()
	url := fmt.Sprintf("%s?token=%s", retriever.Settings.SlackUrl, retriever.Settings.SlackAuthToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return
	}

	var content SlackEmoji
	err = json.Unmarshal(body, &content)
	if err != nil {
		logrus.Error(err)
		return
	}

	var emojis []data.Emo
	for name := range content.Emoji {
		emojis = append(emojis, data.Emo{name})
	}

	retriever.EmoUpdateChannel <- &emojis
	retriever.Metrics.EmoRetrievalLengths.Set(float64(len(emojis)))
	retriever.Metrics.EmoRetrievalDurations.Observe(float64(time.Since(started).Microseconds()))
}

type SlackEmoji struct {
	Ok       bool
	Emoji    map[string]string
	Cache_ts string
}
