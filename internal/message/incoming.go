package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

var (
	PlaneStatusHeader = [3]byte{0x41, 0x49, 0x52}
)

type PlaneStatus struct {
	// International aircraft registration code
	TailNumber string
	// number of engines on the aircraft
	EngineCount uint32
	// engine name
	EngineName string
	// Latitude in degrees
	Latitude float64
	// Longitude in degrees
	Longitude float64
	// Altitude in feet
	Altitude float64
	// Temperature in degrees Fahrenehit
	Temperature float64
}

// Unmarshal a binary byte array into this struct
func (im *PlaneStatus) UnmarshalBinary(data []byte) error {
	// we could add a check for 'data' minimum length
	// using a reader is simpler than keeping track of current read location.
	// if there's a performance hit because of this it can be refactored quite easily
	buffer := bytes.NewBuffer(data)
	// check the header
	var header [3]byte
	if err := binary.Read(buffer, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("reading header: %w", err)
	}
	if header != PlaneStatusHeader {
		return ErrInvalidMessageHeader
	}

	tailNumber, err := readStr(buffer)
	if err != nil {
		return fmt.Errorf("reading tail number: %w", err)
	}
	im.TailNumber = tailNumber

	count, err := readInt(buffer)
	if err != nil {
		return fmt.Errorf("reading engine count: %w", err)
	}
	im.EngineCount = count

	name, err := readStr(buffer)
	if err != nil {
		return fmt.Errorf("reading engine name: %w", err)
	}
	im.EngineName = name

	lat, err := readFloat(buffer)
	if err != nil {
		return fmt.Errorf("reading latitude: %w", err)
	}
	im.Latitude = lat

	long, err := readFloat(buffer)
	if err != nil {
		return fmt.Errorf("reading longitude: %w", err)
	}
	im.Longitude = long

	alt, err := readFloat(buffer)
	if err != nil {
		return fmt.Errorf("reading altitude: %w", err)
	}
	im.Altitude = alt

	temp, err := readFloat(buffer)
	if err != nil {
		return fmt.Errorf("reading temperature: %w", err)
	}
	im.Temperature = temp
	return nil
}

func readFloat(buffer io.Reader) (float64, error) {
	var read [8]byte
	if err := binary.Read(buffer, binary.BigEndian, &read); err != nil {
		return 0.0, err
	}
	bits := binary.BigEndian.Uint64(read[:])
	return math.Float64frombits(bits), nil
}

func readStr(buffer io.Reader) (string, error) {
	var read [4]byte
	if err := binary.Read(buffer, binary.BigEndian, &read); err != nil {
		return "", fmt.Errorf("reading size: %w", err)
	}
	size := binary.BigEndian.Uint32(read[:])
	arr := make([]byte, size)
	if err := binary.Read(buffer, binary.BigEndian, &arr); err != nil {
		return "", err
	}
	return string(arr), nil
}

func readInt(buffer io.Reader) (uint32, error) {
	var read [4]byte
	if err := binary.Read(buffer, binary.BigEndian, &read); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(read[:]), nil
}
