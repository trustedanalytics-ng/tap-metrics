package main

// THERE ARE NO ORGS AT THIS TIME IN TAP

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gocraft/web"
	"gopkg.in/square/go-jose.v1/json"
	"fmt"
)

type serviceContext struct {
	mp MetricsProvider
}

type requestContext struct {
	token         string
	isAdmin       bool
	organizations []string
	platform      bool
}

func handleError(msg string, err error, rw web.ResponseWriter, status int) {
	errMsg := msg
	if err != nil {
		errMsg = errMsg + " : " + err.Error()
	}
	log.Print(errMsg)
	http.Error(rw, errMsg, status)
}

func jsonResponse(rw web.ResponseWriter, response interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(response)
	if err != nil {
		handleError("Error writting response to JSON", err, rw, http.StatusInternalServerError)
	}
}

func (sc *serviceContext) Health(c *requestContext, rw web.ResponseWriter, req *web.Request) {
	// TODO
	rw.Write([]byte("OK"))
}

func extractTime(req *web.Request, key string) (*time.Time, error) {
	// IS Needed?
	t := req.PathParams[key]
	if t == "" {
		return nil, nil
	}
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp := time.Unix(i, 0)
	return &timestamp, nil
}

func getRanges(req *web.Request) (*time.Time, *time.Time) {
	// IS Needed?
	from, _ := extractTime(req, "from")
	to, _ := extractTime(req, "to")
	return from, to
}

func (sc *serviceContext) PlatformMetrics(c *requestContext, rw web.ResponseWriter, req *web.Request) {
	pm, err := sc.mp.PlatformMetrics()
	if err != nil {
		handleError("Error when retrieving platform metrics: ", err, rw, http.StatusInternalServerError)
		return
	}
	jsonResponse(rw, pm)
}

func (sc *serviceContext) OrgMetrics(c *requestContext, rw web.ResponseWriter, req *web.Request) {
	org := req.PathParams["organization"]
	if org != "" {
		handleError("No organization specified", nil, rw, http.StatusBadRequest)
		return
	}
	om, err := sc.mp.OrganizationMetrics(org)
	if err != nil {
		handleError("Error when retrieving platform metrics",
			err, rw, http.StatusInternalServerError)
		return
	}
	jsonResponse(rw, om)
}

func (sc *serviceContext) RawQuery(c *requestContext, rw web.ResponseWriter, req *web.Request) {
	// Consider it as a debug access for now, to be evaluated if needed later on
	query := req.URL.Query().Get("q")
	if query == "" {
		handleError("No query provided", nil, rw, http.StatusBadRequest)
		return
	}
	resp, err := sc.mp.RawQuery(query)
	if err != nil {
		handleError("Error when retrieving metrics for raw query", err, rw, http.StatusNotFound)
		return
	}
	jsonResponse(rw, resp)
}

func (sc *serviceContext) SingleMetric(c *requestContext, rw web.ResponseWriter, req *web.Request) {
	// TODO handle metric with ranges
	// TODO: validate for regexp matching those field (SQL Injection prevention)
	query := req.URL.Query()
	//org := query.Get("organization")
	metric := query.Get("metric")
	from := query.Get("from")
	to := query.Get("to")
	if metric == "" || from == "" || to == "" {
		errMsg := fmt.Sprintf("Not all required parameters were specified: mertic=%s from=%s to=%s",
			metric, from, to)
		handleError(errMsg, nil, rw, http.StatusBadRequest)
		return
	}
	// TODO construct query
	rawMetric, err := sc.mp.SingleMetric("default", []string{"aa", "bb"}, "123", "456")
	if err != nil {
		handleError("Error when retrieving metrics for SingleMetric", err, rw, http.StatusInternalServerError)
		return
	}
	jsonResponse(rw, rawMetric)
}

func (rc *requestContext) authMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	// TODO For now you only need valid TAP UAA token
	// TODO work on this token
	rc.token = "aaa"
	next(rw, r)
}

func setupRouter(mp MetricsProvider) *web.Router {
	sc := serviceContext{mp}
	return web.New(requestContext{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*requestContext).authMiddleware).
		Get("/healthz", sc.Health).
		Get("/api/v1/metrics/platform", sc.PlatformMetrics).
		Get("/api/v1/metrics/organization", sc.OrgMetrics).
		Get("/api/v1/metrics/single", sc.SingleMetric).
		Get("/api/v1/metrics/rawQuery", sc.RawQuery)
}

func setupMetricsProvider() MetricsProvider {
	mp, err := NewInfluxDBMetricsProvider()
	if err != nil {
		log.Fatal("Error while setuping MetricsProvider", err)
	}
	return mp
}

func main() {
	log.Println("Starting Metrics Presenter")
	mp := setupMetricsProvider()
	router := setupRouter(mp)
	err := http.ListenAndServe(":8081", router)
	log.Fatal("Exiting with: ", err)
}
