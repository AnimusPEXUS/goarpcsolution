package goarpcsolution

import (
	"github.com/AnimusPEXUS/gouuidtools"
)

type ARPCTransmissionInfo struct {
	HumanTitle       string
	HumanDescription string
	BufferIds        []*gouuidtools.UUID
}

type ARPCTransmissionI interface{}
