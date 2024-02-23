package handler

import (
	"fmt"
	"io"

	"github.com/broganross/salad-interview/internal/message"
)

// CO2Calculator provides methods for determining a plane's CO2 footprint
type CO2Calculator interface {
	CO2FootPrint(*message.PlaneStatus) error
}

// PlaneStatusHandler knows how to deal with messages formatted for plane status information
type PlaneStatusHandler struct {
	Domain CO2Calculator
}

func (psh *PlaneStatusHandler) ReadConnection(r io.ReadSeeker) (bool, error) {
	m := message.PlaneStatus{}
	data, err := io.ReadAll(r)
	if err != nil {
		r.Seek(0, 0)
		return false, fmt.Errorf("reading: %w", err)
	}
	if err := m.UnmarshalBinary(data); err != nil {
		r.Seek(0, 0)
		return false, fmt.Errorf("unmarshaling PlaneStatus: %w", err)
	}
	fmt.Printf("%+v\n", m)
	// convert to domain struct... probably
	if err := psh.Domain.CO2FootPrint(&m); err != nil {
		r.Seek(0, 0)
		return false, fmt.Errorf("processing fuel usage: %w", err)
	}
	return true, nil
}
