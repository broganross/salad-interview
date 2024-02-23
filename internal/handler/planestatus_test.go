package handler_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/broganross/salad-interview/internal/handler"
	"github.com/broganross/salad-interview/internal/message"
)

type mockDomain struct{}

func (md *mockDomain) CO2FootPrint(m *message.PlaneStatus) error {
	return nil
}

func TestPlaneStatus_ReadConnection(t *testing.T) {
	mocked := mockDomain{}
	tests := []struct {
		name     string
		message  []byte
		expected bool
		err      error
	}{
		{
			"happy-path",
			[]byte{
				0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x06, 0x4E, 0x32, 0x30, 0x39, 0x30, 0x34, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00, 0x07, 0x47, 0x45, 0x6E, 0x78, 0x2D, 0x31, 0x42, 0x40, 0x43, 0x8E, 0xD6,
				0xEB, 0xFF, 0x60, 0x1D, 0xC0, 0x50, 0xD4, 0xC0, 0x91, 0x63, 0x01, 0x65, 0x40, 0xE2, 0x03, 0xF0,
				0x00, 0x00, 0x00, 0x00, 0xC0, 0x4A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A,
			},
			true,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := handler.PlaneStatusHandler{
				Domain: &mocked,
			}
			reader := bytes.NewReader(test.message)
			consumed, err := handler.ReadConnection(reader)
			if !errors.Is(err, test.err) {
				t.Errorf("expected error '%v' got '%v'", test.err, err)
			}
			if consumed != test.expected {
				t.Errorf("expected bool '%v' got '%v'", test.expected, consumed)
			}
		})
	}
}
