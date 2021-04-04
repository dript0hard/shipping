package main

import (
    "bytes"
    "context"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/go-kit/kit/endpoint"
)

// Request, Response modeling for the string service.
type uppercaseRequest struct {
    S string `json:"s"`
}

type uppercaseResponse struct {
    S   string `json:"s"`
    Err string `json:"error,omitempty"`
}

type countRequest struct {
    S string `json:"s"`
}

type countResponse struct {
    L   int `json:"length"`
}

// Endpoint modeling for the string service.

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(uppercaseRequest)
        v, err := svc.Uppercase(req.S)
        if err != nil {
            return uppercaseResponse{S: "", Err: err.Error()}, nil
        }
        return uppercaseResponse{S: v, Err: ""}, nil
    }
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
    return func(_ context.Context, request interface{}) (interface{}, error) {
        req := request.(countRequest)
        v := svc.Count(req.S)
        return countResponse{L: v}, nil
    }
}

func decodeUppercaseRequest(
    _ context.Context, r *http.Request) (interface{}, error) {
    var request uppercaseRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, err
    }
    return request, nil
}

func decodeCountRequest(
    _ context.Context, r *http.Request) (interface{}, error) {
    var request countRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, err
    }
    return request, nil
}

func defaultResponseEncoder(
    _ context.Context, w http.ResponseWriter, response interface{}) error {
    return json.NewEncoder(w).Encode(response)
}

func decodeResponse(_ context.Context, r *http.Response) (interface{}, error){
    var response interface{}
    if err := json.NewDecoder(r.Body).Decode(response); err != nil {
        return nil, err
    }
    return response, nil
}

func encodeRequest(
    _ context.Context, r *http.Request, request interface{}) error {
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(request); err != nil {
        return err
    }

    r.Body = ioutil.NopCloser(&buf)
    return nil
}
