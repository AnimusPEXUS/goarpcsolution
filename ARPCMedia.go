package goarpcsolution

import (
	"errors"
	"net"
)

type ARPCMediaBase struct {
	id string
}

type ARPCMediaBuffer struct {
	ARPCMediaBase
	bufferInfo          func(*ARPCMediaBuffer) any
	readwriteseekcloser ReadWriteSeekCloser
}

// stream isn't required to have underlying buffer.
// in both cases, buffer_size_tracker used to truck size of buffer
// and to calculate next 'index' to be sent to listeners
type ARPCMediaBroadcast struct {
	ARPCMediaBase
	Buffer              *ARPCMediaBuffer
	buffer_size_tracker int64
	listeners           []func(
		stream *ARPCMediaBroadcast,
		// this is essentially stream.Buffer
		// just a shortcut
		buffer *ARPCMediaBuffer,
		// data it self. (end - start) must be rql to len(data)
		data []byte,
		index int64,
	)
}

type ARPCMediaConn struct {
	ARPCMediaBase
	OnAccept func(s *ARPCMediaConn, c net.Conn) error
}

/////////////// Base Functions

func (self *ARPCMediaBase) Id() string {
	return self.id
}

/////////////// Buffer Functions

func (self *ARPCMediaBuffer) GetReadWriteSeekCloser() ReadWriteSeekCloser {
	return self.readwriteseekcloser
}

/////////////// Broadcast Functions

// recommended maximum data size - 1024 bytes
// fixed maximum data size - 2 MiB
func (self *ARPCMediaBroadcast) PushFeed(data []byte) error {
	data_size := int64(len(data))

	// those checks vere checking parameters, as index and end was
	// parameters. checks are left and will be ermoved at some point
	// todo: remove checks?
	if data_size > (2 * 1024 * 1024) {
		return errors.New("bufer is too big. 2MiB maximum")
	}

	index := self.buffer_size_tracker

	if index < 0 {
		return errors.New("invalid 'index' value")
	}

	end := index + data_size

	return nil
}

/////////////// Socket Functions
