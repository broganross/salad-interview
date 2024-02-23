# Technical Project for Salad Technologies
Original problem description is found here: https://saladtech.notion.site/Salad-s-Technical-Interview-2e7e8c19c4284c61975fd3d4eba3bfbe

## Configuration
Configuration is managed with the following environment variables

#### SALAD_MESSAGEROUTER_URL
&ensp; The URL for the message router.

&ensp; **Example**: `data.salad.com:5000`

&ensp; **Default**:

&ensp; **Required**: Yes

#### SALAD_MESSAGEROUTER_TCPTIMEOUT
&ensp; Timeout control for open TCP connections.  If set to 0, no timeout is set

&ensp; **Example**: 30s

&ensp; **Default**: 1m

&ensp; **Required**: No

#### SALAD_MESSAGEROUTER_MAXRETRY
&ensp; Maximum number of times to try opening a connection to the server before exiting.
&ensp; -1 means there's no maximum.

&ensp; **Example**: 10

&ensp; **Default**: 5

&ensp; **Required**: No

#### SALAD_MESSAGEROUTER_RETRYSLEEP 
&ensp; The starting duration to wait between retrying TCP connections.  This value will increase exponentially for each retry.

&ensp; **Example**: 3s

&ensp; **Default**: 1s

&ensp; **Required**: No

#### SALAD_LOG_LEVEL
&ensp; Level for the logger
&ensp; Available log level values are as follows:
* Debug: 0
* Info: 1
* Warn: 2
* Error: 3
* Fatal: 4
* Panic: 5

&ensp; **Example**: 2

&ensp; **Default**: 1

&ensp; **Required**: No

## Execute Unit Tests

    go test ./...

## Build

    go build ./cmd/processor/main.go


## Run

    go run ./cmd/processor/main.go


## Notes

1. The problem says the float64 fields have a size of 4 bytes.  This obviously doesn't work, but the example does contain enough bytes.
1. Resolving the name 'data.salad.com' fails.  An nslookup makes it seem like it's an example.  Therefore I made a simple TCP server which just gives a valid response.  It can be run using: ```go run ./cmd/server/main.go```.  Then set the SALAD_MESSAGEROUTER_URL environment variable to `127.0.0.1:5000`.