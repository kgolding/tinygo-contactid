package main

import (
	"runtime"
	"strings"
	"time"

	"tinygo.org/x/drivers"
)

var (
	OK = []byte{'O', 'K', 0x13, 0x0a}
)

const (
	minimumLength     = 19
	startingDelimiter = 'A'
	endDelimiter      = 0x0d
	bufferSize        = 8
)

// Device wraps a connection to a device.
type Device struct {
	buffer   []byte
	bufIdx   int
	bufLen   int
	sentence strings.Builder
	uart     drivers.UART
}

// NewContactIdUart creates a new Contact ID connection over a pre-configured UART connection
func NewContactIdUart(uart drivers.UART) Device {
	return Device{
		uart:     uart,
		buffer:   make([]byte, bufferSize),
		bufIdx:   bufferSize,
		sentence: strings.Builder{},
	}
}

// NextMessage returns the next valid message (blocking)
func (d *Device) NextMessage() (msg *CIDMessage, err error) {
	s, err := d.NextSentence()
	if err != nil {
		return nil, err
	}
	return parseCID(s)
}

// NextSentence returns the next valid sentence from the device, and
// replies "OK" (blocking)
func (d *Device) NextSentence() (sentence string, err error) {
	sentence = d.readNextSentence()
	if err = validSentence(sentence); err != nil {
		return "", err
	}
	d.uart.Write(OK)
	return sentence, nil
}

// readNextSentence returns the next sentence from the device.
func (d *Device) readNextSentence() (sentence string) {
	d.sentence.Reset()
	var b byte = ' '

	// Find the starting delimiter
	for b != startingDelimiter {
		b = d.readNextByte()
	}

	for b != endDelimiter {
		d.sentence.WriteByte(b)
		b = d.readNextByte()
	}
	d.sentence.WriteByte(b)

	sentence = d.sentence.String()
	return sentence
}

func (d *Device) readNextByte() (b byte) {
	d.bufIdx += 1
	if d.bufIdx >= d.bufLen {
		d.fillBuffer()
	}
	return d.buffer[d.bufIdx]
}

func (d *Device) fillBuffer() {
	for d.uart.Buffered() == 0 {
		time.Sleep(time.Millisecond * 100)
		runtime.Gosched()
	}
	d.bufLen, _ = d.uart.Read(d.buffer[0:bufferSize])
	d.bufIdx = 0
	// println("<<< ", string(d.buffer[:d.bufLen]))
}

// WriteBytes sends data to the device
func (d *Device) WriteBytes(bytes []byte) {
	d.uart.Write(bytes)
}

// validSentence checks if a sentence has been received
func validSentence(sentence string) error {
	if len(sentence) < minimumLength || sentence[0] != startingDelimiter {
		return errInvalidSentence
	}
	return nil
}
