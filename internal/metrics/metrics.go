package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	EmoRetrievalDurations prometheus.Histogram
	EmoRetrievalLengths   prometheus.Gauge
	EmoRetrievalCount     prometheus.Counter

	RandomEmoRequestCount     *prometheus.CounterVec
	RandomEmoRequestDurations *prometheus.HistogramVec
}

func NewMetrics() Metrics {
	namespace := "random_slack_emo"
	emojiRetrievalSubsystem := "emoji_retrieval"
	randomEmojiSubsystem := "random_emoji_requests"

	metrics := Metrics{
		EmoRetrievalDurations: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: emojiRetrievalSubsystem,
			Name:      "emoji_retrieval_duration",
			Help:      "Tracks how long it takes to retrieve the emoticons from Slack",
			Buckets:   prometheus.DefBuckets,
		}),
		EmoRetrievalLengths: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: emojiRetrievalSubsystem,
			Name:      "emoji_lengths",
			Help:      "The number of emojis returned by Slack",
		}),
		EmoRetrievalCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: emojiRetrievalSubsystem,
			Name:      "emoji_retrieval_count",
			Help:      "The number of requests to retrieve emojis from Slack",
		}),
		RandomEmoRequestCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: randomEmojiSubsystem,
			Name:      "random_emoji_request_count",
			Help:      "The number of requests to retrieve random a random emoji",
		}, []string{"code", "method"}),
		RandomEmoRequestDurations: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: randomEmojiSubsystem,
			Name:      "random_emoji_request_duration",
			Help:      "Tracks how long it takes to retrieve a random emoji",
			Buckets:   prometheus.DefBuckets,
		}, []string{"code", "method"}),
	}

	err := prometheus.Register(metrics.EmoRetrievalLengths)
	err = prometheus.Register(metrics.EmoRetrievalCount)
	err = prometheus.Register(metrics.EmoRetrievalDurations)
	err = prometheus.Register(metrics.RandomEmoRequestDurations)
	err = prometheus.Register(metrics.RandomEmoRequestCount)

	println(err)

	return metrics
}
