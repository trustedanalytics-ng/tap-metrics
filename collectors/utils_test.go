package utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	MEASUREMENT            = "test_measurement"
	TIMESTAMP              = 1234567890123
	BUFF_SIZE              = 10
	CRON_PATTERN           = "* * * * *"
	ITERATIONS             = 2
	EXPECTED_SERIALIZATION = "test_measurement,t1=t1val,t2=t2val f1=1,f2=3.14,f3=someStr 1234567890123"
)

type TestCollector struct{}

func (*TestCollector) CollectMetrics() ([]*Metrics, error) {
	testMetric := Metrics{
		MEASUREMENT,
		map[string]interface{}{
			"f1": 1,
			"f2": 3.14,
			"f3": "someStr",
		},
		map[string]string{
			"t1": "t1val",
			"t2": "t2val",
		},
		TIMESTAMP,
	}
	return []*Metrics{&testMetric}, nil
}

func TestE2EUtils(t *testing.T) {
	metricsChan := make(chan string, BUFF_SIZE)
	collectorAddress := startCollectorServiceMock(metricsChan)

	collectingEngine := prepareCollectingEngine(collectorAddress)
	collectingEngine.Start()

	time.Sleep((1000 * ITERATIONS) * time.Millisecond)
	collectingEngine.BlockingStop()

	validateAllMetrics(t, metricsChan)
}

func validateAllMetrics(t *testing.T, metricsChan chan string) {
	count := 0
	for metric := range metricsChan {
		if EXPECTED_SERIALIZATION != metric {
			t.Errorf("Unexpected metric: '%s'", metric)
		}
		count++
	}
	if count != ITERATIONS {
		t.Errorf("Wrong number of metrics. Got %d, expected %d", count, ITERATIONS)
	}
}

func startCollectorServiceMock(metricsChan chan string) string {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("Failed creating tcp listener", err)
	}
	go handleSingleConnection(metricsChan, ln)
	return ln.Addr().String()
}

func handleSingleConnection(metricsChan chan string, listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Failed accepting tcp connection", err)
	}
	defer conn.Close()

	buff := make([]byte, 4096)
	str := ""
	for {
		n, err := conn.Read(buff)
		fmt.Printf("TCP server recieved len: %d err: %v\n", n, err)
		if n != 0 {
			str += string(buff[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading the data from socet")
		}
	}
	for _, metric := range strings.Split(str, "\n") {
		if len(metric) > 0 {
			metricsChan <- metric
		}
	}
	close(metricsChan)
}

func prepareCollectingEngine(collectorAddress string) *CollectingEngine {
	collectingEngine, err := NewCollectingEngine(collectorAddress, BUFF_SIZE)
	if err != nil {
		log.Fatal("Error creating Collecting Engine", err)
	}
	var metricsCollector MetricsCollector = &TestCollector{}
	collectingEngine.AddMetricsCollector(CRON_PATTERN, &metricsCollector)
	return collectingEngine
}

func TestSerialization(t *testing.T) {
	metric, _ := (&TestCollector{}).CollectMetrics()
	serializedMetric := serializeMetric(metric[0])
	if serializedMetric != EXPECTED_SERIALIZATION {
		t.Errorf("Wrong serialization of metrics: ", serializedMetric)
	}
}
