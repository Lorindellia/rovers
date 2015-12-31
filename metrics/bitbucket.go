package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subsystemBitbucket = subsystem + "_bitbucket"

	BitbucketProcessed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: subsystemBitbucket,
			Name:      "processed",
			Help:      "Number of Bitbucket repositories processed.",
		},
	)
	BitbucketFailed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: subsystemBitbucket,
			Name:      "failed",
			Help:      "Number of Bitbucket repositories failed.",
		},
		[]string{"reason"},
	)
	BitbucketRequested = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: subsystemBitbucket,
			Name:      "requested",
			Help:      "number of requests made to Bitbucket API.",
		},
	)
	BitbucketRequestDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemBitbucket,
			Name:      "request_duration_microseconds",
			Help:      "The Bitbucket API request latency in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(BitbucketProcessed)
	prometheus.MustRegister(BitbucketFailed)
	prometheus.MustRegister(BitbucketRequested)
	prometheus.MustRegister(BitbucketRequestDur)
}