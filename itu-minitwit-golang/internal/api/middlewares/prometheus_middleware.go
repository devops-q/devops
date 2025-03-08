package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
)

var (
	CpuGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percentage",
			Help: "CPU usage in percentage",
		},
		[]string{"cpu", "host"},
	)

	ResponseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_responses_count",
			Help: "Count of HTTP responses sent",
		},
		[]string{"status_code"},
	)

	ReqDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "Request duration in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "route"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(CpuGauge)
	prometheus.MustRegister(ResponseCounter)
	prometheus.MustRegister(ReqDurationSummary)
}

// UpdateCPUMetrics sets the CPU metrics
func UpdateCPUMetrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get CPU usage
		percent, err := cpu.Percent(0, false)
		if err != nil {
			log.Printf("Error getting CPU usage: %v", err)
		}

		// Set the CPU usage metric
		CpuGauge.WithLabelValues("cpu0", "host1").Set(percent[0])

		// Call the next handler
		ctx.Next()
	}
}

// TrackRequestDuration and count responses
func TrackRequestDuration() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Record the start time of the request
		start := time.Now()

		// Process the request
		ctx.Next()

		// Calculate request duration and observe it in Prometheus
		duration := time.Since(start).Seconds()
		ReqDurationSummary.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(duration)

		// Increment the response counter with status code
		statusCode := ctx.Writer.Status()
		ResponseCounter.WithLabelValues(string(rune(statusCode))).Inc()
	}
}
