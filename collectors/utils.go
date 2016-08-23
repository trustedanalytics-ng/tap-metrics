package utils

// I'm sorry for dummy commens but my IDE complains about them on every save
// and I can't turn it off

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/robfig/cron"
)

// Metrics represent data put into InfluxDB
type Metrics struct {
	Measurement string
	Fields      map[string]interface{}
	Tags        map[string]string
	Timestamp   int64
}

// When you want to write specyfic metric collector this is interface that you need to implement
type MetricsCollector interface {
	// Can return Metrics and error at the same time
	CollectMetrics() ([]*Metrics, error) // Maybe list of errors for specyfic metrics?
}

// Main component that you use to register your MetricCollectors and schedule them
type CollectingEngine struct {
	cron          *cron.Cron
	metricsChan   chan *Metrics
	done          chan bool
	wgCollectors  sync.WaitGroup
	collectorConn net.Conn
}

func NewCollectingEngine(mainCollectorAddress string, bufferSize int) (*CollectingEngine, error) {
	cron := cron.New()

	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", mainCollectorAddress)
	if err != nil {
		log.Printf("Unable to connect to %s", mainCollectorAddress)
		return nil, err
	}

	metricsChan := make(chan *Metrics, bufferSize)
	var wg sync.WaitGroup
	done := make(chan bool)

	return &CollectingEngine{cron, metricsChan, done, wg, conn}, nil
}

// **** Metrics collection ****

func (ce *CollectingEngine) collectMetricsToChan(mc MetricsCollector) {
	timestamp := time.Now().UnixNano()
	metrics, err := mc.CollectMetrics()
	if err != nil {
		log.Printf("Error while fetching metrics %v : %v, but got %d\n",
			mc, err, len(metrics))
	}
	for i := range metrics {
		if metrics[i].Timestamp == 0 {
			metrics[i].Timestamp = timestamp
		}
		ce.metricsChan <- metrics[i]
	}
	ce.wgCollectors.Done()
}

func (ce *CollectingEngine) collectMetrics(mcs []*MetricsCollector) {
	ce.wgCollectors.Add(len(mcs))
	for i := range mcs {
		go ce.collectMetricsToChan(*mcs[i])
	}
}

// **** Metrics sending ****

func (ce *CollectingEngine) sendingLoop() {
	for metric := range ce.metricsChan {
		ce.sendMetric(metric)
	}
	ce.done <- true
}

func sortedKeys(m map[string]interface{}) []string {
	sortedKeys := make([]string, len(m))
	i := 0
	for k := range m {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}

func serializeMap(m map[string]interface{}) string {
	sortedKeys := sortedKeys(m)
	serialized := make([]string, 0, len(m))
	for _, key := range sortedKeys {
		kvSerialized := fmt.Sprintf("%v=%v", key, m[key])
		serialized = append(serialized, kvSerialized)
	}
	return strings.Join(serialized, ",")
}

func generalizeMap(m map[string]string) map[string]interface{} {
	genMap := make(map[string]interface{}, len(m))
	for key := range m {
		genMap[key] = m[key]
	}
	return genMap
}

func serializeMetric(metric *Metrics) string {
	key := metric.Measurement
	tags := serializeMap(generalizeMap(metric.Tags))
	fields := serializeMap(metric.Fields)
	if len(tags) > 0 {
		key = key + "," + tags
	}
	return fmt.Sprintf("%s %s %v", key, fields, metric.Timestamp)
}

func (ce *CollectingEngine) sendMetric(metric *Metrics) {
	serializedMetric := serializeMetric(metric)
	fmt.Fprintf(ce.collectorConn, "%s\n", serializedMetric)
}

// **** Public interface ****

func (ce *CollectingEngine) Start() {
	go ce.sendingLoop()
	ce.cron.Start()
}

func (ce *CollectingEngine) BlockingStop() {
	ce.cron.Stop()
	ce.wgCollectors.Wait()
	close(ce.metricsChan)
	<-ce.done
	ce.collectorConn.Close()
}

func (ce *CollectingEngine) AddMetricsCollector(schedule string, mc *MetricsCollector) {
	ce.AddMetricsCollectors(schedule, []*MetricsCollector{mc})
}

func (ce *CollectingEngine) AddMetricsCollectors(schedule string, mcs []*MetricsCollector) {
	ce.cron.AddFunc(schedule, func() { ce.collectMetrics(mcs) }) // this already will be in goroutine
}

func (ce *CollectingEngine) StartAndBlockTillSigInt() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ce.Start()
	log.Println("Started CollectingEngine")
	<-c
	log.Println("Got SIGINT, shutting down metrics collecting")
	ce.BlockingStop()
}
