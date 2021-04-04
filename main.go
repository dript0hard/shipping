package main

import (
    "os"
    "net/http"
    "github.com/go-kit/kit/log"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
    httptransport "github.com/go-kit/kit/transport/http"
)

// Trasport layer with entrypoint.
func main() {
    logger := log.NewLogfmtLogger(os.Stderr)

    fieldKeys := []string{"method", "error"}

    requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name: "request_count",
        Help: "Number of requests recieved.",
    }, fieldKeys)

    requestLarency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name: "request_latency_microseconds",
        Help: "Total duration of requests in microseconds",
    }, fieldKeys)

    countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
        Namespace: "my_group",
        Subsystem: "string_service",
        Name: "count_result",
        Help: "The result of each count method.",
    }, []string{})

    // Setup middleware
    var svc StringService
    svc = stringService{}
    svc = loggingMiddleware{logger: logger, next: svc}
    svc = instrumentMiddleware{
        requestCount: requestCount,
        requestLarency: requestLarency,
        countResult: countResult,
        next: svc,
    }

    // Setup Serve Https
    uppercaseHandler := httptransport.NewServer(
        makeUppercaseEndpoint(svc),
        decodeUppercaseRequest,
        defaultResponseEncoder,
    )

    countHandler := httptransport.NewServer(
        makeCountEndpoint(svc),
        decodeCountRequest,
        defaultResponseEncoder,
    )

    // Setup handlers
    http.Handle("/uppercase", uppercaseHandler)
    http.Handle("/count", countHandler)
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
