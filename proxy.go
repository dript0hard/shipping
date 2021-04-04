package main

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "net/url"

    "golang.org/x/time/rate"

    "github.com/sony/gobreaker"

    "github.com/go-kit/kit/sd"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/sd/lb"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/ratelimit"
    "github.com/go-kit/kit/circuitbreaker"
    httptransport "github.com/go-kit/kit/transport/http"
)

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
    next      StringService
    uppercase endpoint.Endpoint
}

func (mw proxymw) Uppercase(s string) (string, error) {
    response, err := mw.uppercase(nil, uppercaseRequest{S: s})
    if err != nil {
        return "", err
    }

    resp := response.(uppercaseResponse)
    if resp.Err != "" {
        return resp.S, errors.New(resp.Err)
    }

    return resp.S, nil
}

func (mw proxymw) Count(s string) int {
    return mw.next.Count(s)
}

func proxyingMiddleware(instances string, logger log.Logger) ServiceMiddleware {

    // If instances is empty, don't proxy the requests.
    if instances == "" {
        logger.Log("proxy_to", "none")
        return func(next StringService) StringService { return next }
    }

    // Set params for our client.
    var (
        qps         = 100
        maxAttempes = 3
        maxTime     = 250 * time.Millisecond
    )

    // Construct and Endpoint for each instance in the list, add it to a
    // fixed set of endpoints. In a real service, rather than doing this
    // by hand, you'd propably use package sd's support for your service
    // discovery system.
    var (
        instanceList = split(instances)
        endpointer sd.FixedEndpointer
    )
    logger.Log("proxy_to", fmt.Sprint(instanceList))

    for _, instance := range instanceList {
        var e endpoint.Endpoint
        e = makeUppercaseProxy(instance)

        // Circuit breaking for failing request.
        e = circuitbreaker.Gobreaker(
            gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)

        // Rate limiting for 100 querys per second.
        e = ratelimit.NewErroringLimiter(
            rate.NewLimiter(rate.Every(time.Second), qps))(e)

        // Add endpoint to enporinter list.
        endpointer = append(endpointer, e)
    }

    // Now build a single, retrying, load-balancing, enpoint out of all of
    // those individual enpoints.
    balancer := lb.NewRoundRobin(endpointer)
    retry    := lb.Retry(maxAttempes, maxTime, balancer)

    return func(next StringService) StringService {
        return proxymw{ next: next, uppercase: retry }
    }
}

func makeUppercaseProxy(proxyUrl string) endpoint.Endpoint {
    return httptransport.NewClient(
        "GET",
        mustParse(proxyUrl),
        encodeRequest,
        decodeResponse,
    ).Endpoint()
}

func mustParse(proxyUrl string) *url.URL {
    u, err := url.Parse(proxyUrl)
    if err != nil {
        panic(err)
    }
    return u
}

func split(instances string) []string {
    a := strings.Split(instances, ",")
    for i := range a {
        a[i] = strings.TrimSpace(a[i])
    }
    return a
}
