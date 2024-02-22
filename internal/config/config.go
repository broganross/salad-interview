package config

import "time"

type Config struct {
	MessageRouter MessageRouter
	Log           Log
}

// The Message router is the service which dispatches messages to our processor
type MessageRouter struct {
	URL       string        `required:"true"`
	TCPTimout time.Duration `default:"3m"`
}

// Logger settings
type Log struct {
	Level int `default:"1"`
}
