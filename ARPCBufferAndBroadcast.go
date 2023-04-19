package goarpcsolution

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ARPCBufferMode uint8

const (
	ARPCBufferModeInvalid ARPCBufferMode = iota

	// if buffer items are byte slices and passed throug calls
	// as byte arrays
	ARPCBufferModeBinary

	// if buffer items are objects and passed through calls
	// as JS objects
	ARPCBufferModeObject
)

type ARPCBufferInfo struct {
	Id               uuid.UUID
	HumanTitle       string
	HumanDescription string

	// todo: todo. this part is in fog for now
	Mode            ARPCBufferMode
	Finished        bool
	OnGoing         bool
	TechDescription any
}

// if buffer mode is 'binary', ItemId must be of form 'x:y'
// like in Go's slice syntax. so x - is index of this items's first byte in
// corresponding buffer, and y is index of last byte in buffer.
// value, in this case, must be []byte type
//
// in 'object mode' , ItemId should be simple integer, corresponding to
// index of this item in buffer. value, in this case, can be any type, which
// can be marshalled via Go's json module
type ARPCBufferItem struct {
	BufferId uuid.UUID
	ItemId   string
	ItemTime time.Time

	Value any
}

type ARPCBroadcastInfo struct {
	HumanTitle       string
	HumanDescription string
	BufferIds        []*uuid.UUID
}

type ARPCBuffer interface {
	GetInfo() *ARPCBufferInfo
	ItemCount() int
	// if not found - it's not error and 2nd result is false
	GetItem(id string) (*ARPCBufferItem, bool, error)
}
