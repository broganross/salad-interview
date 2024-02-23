package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/broganross/salad-interview/internal/config"
	"github.com/rs/zerolog/log"
)

type Handler interface {
	// Reads the message and returns whether it consumed the messasge, and possibly an error
	ReadConnection(io.ReadSeeker) (bool, error)
}

// Listener connects to the message server, and continuously polls for messages
type Listener struct {
	// URL to connect to
	URL string
	// Sets a deadline for the TCP connection
	Timeout time.Duration
	// Maximum number of times to retry a connection before reporting an issue
	// if set to -1 retries will happen infinitely
	MaxRetry int
	// Starting amount of time to sleep between retrying the connection
	// this value doubles each time a retry if executed
	RetrySleep time.Duration
	// Handlers to pass the message into
	handlers []Handler
	ctx      context.Context
}

func NewListener(ctx context.Context, conf *config.Config) *Listener {
	return &Listener{
		ctx:        ctx,
		URL:        conf.MessageRouter.URL,
		Timeout:    conf.MessageRouter.TCPTimout,
		MaxRetry:   conf.MessageRouter.MaxRetry,
		RetrySleep: conf.MessageRouter.RetrySleep,
	}
}

func (l *Listener) makeConnection() (net.Conn, error) {
	log.Debug().Str("URL", l.URL).Msg("making connection")
	retryCount := l.MaxRetry
	sleep := l.RetrySleep
	var conn net.Conn
	var err error
	// create a connection with expotential backoff, if configured to
	for {
		conn, err = net.Dial("tcp", l.URL)
		if err != nil {
			if retryCount == 0 {
				err = ErrMaxRetry
				break
			}
			log.Error().Err(err).Msg("dialing TCP")
			time.Sleep(sleep)
			sleep *= 2
			retryCount--
			log.Debug().Dur("sleep", sleep).Int("retryCount", retryCount).Msg("retrying connection")
		} else {
			break
		}
	}
	// TODO: If it's important for the listener to really just wait for a message,
	// we would remove this functionality
	if conn != nil && l.Timeout.Microseconds() > 0 {
		if err := conn.SetDeadline(time.Now().Add(l.Timeout)); err != nil {
			return conn, err
		}
	}
	return conn, err
}

// AddHandler adds the given handler to the list which will run on a message
func (l *Listener) AddHandler(h Handler) error {
	// This could probably use some validation
	l.handlers = append(l.handlers, h)
	return nil
}

// Starts making connections and waiting for the server to send messages
func (l *Listener) Listen() error {
	log.Info().Msg("listener starting")
	var returnErr error
	for {
		conn, err := l.makeConnection()
		if errors.Is(err, ErrMaxRetry) {
			// TODO: Determine how we actually want to handle this.
			returnErr = err
			break
			// or
			// time.Sleep(1 * time.Second)
			// continue
		} else if err != nil {
			// TODO: Do we really want to just keep trying infinitely?
			time.Sleep(1 * time.Second)
			continue
		}
		// listen on the connection until it writes something
		readerChan := make(chan io.ReadSeeker, 1)
		errChan := make(chan error, 1)
		go func() {
			// TODO: it might be better to wrap the net.Conn in another struct instead
			// of reading it all right here
			raw, err := io.ReadAll(conn)
			defer conn.Close()
			if err != nil {
				errChan <- err
			} else {
				readerChan <- bytes.NewReader(raw)
			}
			close(readerChan)
			close(errChan)
		}()

		select {
		case <-l.ctx.Done():
			log.Debug().Msg("closing connection")
			conn.Close()
			returnErr = l.ctx.Err()
		case reader := <-readerChan:
			// we got a message, pass it into all the handlers
			go func(r io.ReadSeeker) {
				for _, handler := range l.handlers {
					consumed, err := handler.ReadConnection(reader)
					if err != nil {
						log.Error().Err(err).Msg("handler reading connection")
					}
					if consumed {
						break
					}
				}
			}(reader)
		case e := <-errChan:
			log.Error().Err(e).Msg("reading connection")
			// TODO: sleep??
		}
		if returnErr != nil {
			break
		}
	}
	return returnErr
}
