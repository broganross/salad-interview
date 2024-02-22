# Technical Project for Salad Technologies
Original problem description is found here: https://saladtech.notion.site/Salad-s-Technical-Interview-2e7e8c19c4284c61975fd3d4eba3bfbe

## Configuration
Configuration is managed with environment variables

    # The URL for the message router
    SALAD_MESSAGEROUTER_URL=data.salad.com:5000
    # Level for the logger
    SALAD_LOG_LEVEL=1

Available log level values are as follows:

* Debug: 0
* Info: 1
* Warn: 2
* Error: 3
* Fatal: 4
* Panic: 5


## Execute Unit Tests

    go test ./...

## Build

    go build ./cmd/processor/main.go


## Run

    go run ./cmd/processor/main.go


## Notes

1. The problem says the float64 fields have a size of 4 bytes.  This obviously wouldn't work, and it turns out the example does contain enough bytes for float64s.