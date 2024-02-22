package listener

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/broganross/salad-interview/internal/config"
	"github.com/broganross/salad-interview/internal/message"
	"github.com/rs/zerolog/log"
)

type Listener struct {
	ctx context.Context
	// URL to connect to
	URL string
	// Sets a deadline for the TCP connection
	Timeout time.Duration
}

func New(ctx context.Context, conf *config.Config) *Listener {
	return &Listener{
		ctx:     ctx,
		URL:     conf.MessageRouter.URL,
		Timeout: conf.MessageRouter.TCPTimout,
	}
}

func (l *Listener) Start() error {
	log.Info().Msg("listener starting")
	go func() {
		retryCount := 0
		sleep := 1 * time.Second
		for {
			// build a connection
			log.Debug().Msg("opening connection")
			conn, err := net.Dial("tcp", l.URL)
			if err != nil {
				log.Error().Str("URL", l.URL).Err(err).Msg("dialing TCP")
				// retry with rollback
				time.Sleep(sleep)
				sleep *= 2
				retryCount++
				log.Debug().Int("sleep", int(sleep)).Int("retryCount", retryCount).Msg("retrying connection")
				continue
			}
			defer conn.Close()
			// reset retry values
			retryCount = 0
			sleep = 1 * time.Second

			if err := conn.SetDeadline(time.Now().Add(l.Timeout)); err != nil {
				log.Error().Err(err).Msg("setting connection deadline")
				continue
			}

			select {
			case <-l.ctx.Done():
				conn.Close()
				log.Error().Err(l.ctx.Err()).Msg("interrupt triggered")
			default:
				raw := make([]byte, 1024)
				// var raw []byte
				raw, err := io.ReadAll(conn)
				// Need to add handling here
				if err != nil {
					log.Error().Err(err).Msg("reading from server")
					continue
				}
				m := message.PlaneStatus{}
				if err := m.UnmarshalBinary(raw); err != nil {
					log.Error().Err(err).Msg("unmarshaling message")
					continue
				}
				fmt.Println(m)
			}
		}
	}()
	return nil
}
