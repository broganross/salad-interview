package middleware

import (
	"io"

	"github.com/rs/zerolog/log"
)

// basic example of a middleware
type Log struct{}

func (l *Log) ReadConnection(r io.ReadSeeker) (bool, error) {
	// Setup logging associated to the message
	// this is a reason why we would want to wrap the message
	// we could pass a context with the message, and use it to
	// make a logger specific to this message
	log.Debug().Msg("middleware proc")
	return false, nil
}
