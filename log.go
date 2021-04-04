package main

import (
    "time"
    "github.com/go-kit/kit/log"
)


type loggingMiddleware struct {
    logger log.Logger
    next   StringService
}

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
    begin := time.Now()
    output, err = mw.next.Uppercase(s)
    mw.logger.Log(
        "time", begin.String(),
        "method", "uppercase",
        "input", s,
        "error", err,
        "took", time.Since(begin),
    )
    return
}

func (mw loggingMiddleware) Count(s string) (n int) {
    begin := time.Now()
    n = mw.next.Count(s)
    mw.logger.Log(
        "time", begin.String(),
        "method", "count",
        "length", n,
        "took", time.Since(begin),
    )
    return
}

