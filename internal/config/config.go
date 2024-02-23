package config

import "time"

// Config uses the "github.com/kelseyhightower/envconfig" library to
// parse env vars to this configuration
type Config struct {
	MessageRouter MessageRouter
	Log           Log
}

// The Message router is the service which dispatches messages to our processor
type MessageRouter struct {
	// SALAD_MESSAGEROUTER_URL
	URL string `required:"true"`
	// SALAD_MESSAGEROUTER_TCPTIMEOUT
	TCPTimout time.Duration `default:"1m"`
	// SALAD_MESSAGEROUTER_MAXRETRY
	MaxRetry int `default:"5"`
	// SALAD_MESSAGEROUTER_RETRYSLEEP
	RetrySleep time.Duration `default:"1s"`
}

// Logger settings
type Log struct {
	// SALAD_LOG_LEVEL
	Level int `default:"1"`
}
