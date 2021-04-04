package main

import(
    "errors"
    "strings"
)

// Api of this service as an interface.
type StringService interface {
    // Returns uppercase string of the input string.
    Uppercase(string) (string, error)
    // Returns the length of the string.
    Count(string) int
}

// Implementation of the StringService interface for this service.
type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
    if s == "" {
        return "", ErrEmpty
    }
    return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
    return len(s)
}

var ErrEmpty = errors.New("Empty string.")

type ServiceMiddleware func(StringService) StringService
