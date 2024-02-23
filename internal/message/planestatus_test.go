package message_test

import (
	"errors"
	"io"
	"testing"

	"github.com/broganross/salad-interview/internal/message"
)

func TestIncomingMessage_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		err      error
		expected message.PlaneStatus
	}{
		{
			"basic",
			[]byte{
				0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x06, 0x4E, 0x32, 0x30, 0x39, 0x30, 0x34, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00, 0x07, 0x47, 0x45, 0x6E, 0x78, 0x2D, 0x31, 0x42, 0x40, 0x43, 0x8E, 0xD6,
				0xEB, 0xFF, 0x60, 0x1D, 0xC0, 0x50, 0xD4, 0xC0, 0x91, 0x63, 0x01, 0x65, 0x40, 0xE2, 0x03, 0xF0,
				0x00, 0x00, 0x00, 0x00, 0xC0, 0x4A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A,
			},
			nil,
			message.PlaneStatus{
				TailNumber:  "N20904",
				EngineCount: 2,
				EngineName:  "GEnx-1B",
				Latitude:    39.11593389482025,
				Longitude:   -67.32425341289998,
				Altitude:    36895.5,
				Temperature: -53.2,
			},
		},
		{
			"empty-buffer",
			[]byte{},
			io.EOF,
			message.PlaneStatus{},
		},
		{
			"wrong-header",
			[]byte{0x42, 0x49, 0x52},
			message.ErrInvalidMessageHeader,
			message.PlaneStatus{},
		},
		{
			"bad-tail-number",
			[]byte{0x41, 0x49, 0x52, 0x00, 0x00, 0x00},
			io.ErrUnexpectedEOF,
			message.PlaneStatus{},
		},
		// TODO: extend these unit tests
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := message.PlaneStatus{}
			err := m.UnmarshalBinary(test.input)
			if !errors.Is(err, test.err) {
				t.Errorf("expected error '%v' got '%v'", test.err, err)
				return
			}
			if test.expected != m {
				t.Errorf("expected parsed message '%v' got '%v'", test.expected, m)
			}
		})
	}
}
