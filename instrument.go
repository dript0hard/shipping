package main

import (
    "fmt"
    "time"
    "github.com/go-kit/kit/metrics"
)

type instrumentMiddleware struct {
    requestCount   metrics.Counter
    requestLarency metrics.Histogram
    countResult    metrics.Histogram
    next           StringService
}

func (mw instrumentMiddleware) Uppercase(s string) (output string, err error){
    begin := time.Now()
    output, err = mw.next.Uppercase(s)
    lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
    mw.requestCount.With(lvs...).Add(1)
    mw.requestLarency.With(lvs...).Observe(time.Since(begin).Seconds())
    return
}

func (mw instrumentMiddleware) Count(s string) (n int) {
    begin := time.Now()
    n = mw.next.Count(s)
    lvs := []string{"method", "count", "error", "false"}
    mw.requestCount.With(lvs...).Add(1)
    mw.requestLarency.With(lvs...).Observe(time.Since(begin).Seconds())
    mw.countResult.Observe(float64(n))
    return
}
