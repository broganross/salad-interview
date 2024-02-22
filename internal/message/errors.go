package message

import "errors"

var (
	ErrInvalidMessageHeader = errors.New("invalid message header")
	ErrMessageTooShort      = errors.New("message not long enough")
)
