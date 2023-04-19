package goarpcsolution

import (
	"errors"
	"runtime"
	"strings"
	"sync"

	"github.com/AnimusPEXUS/goinmemfile"
	uuid "github.com/satori/go.uuid"
)

type ARPCMediaObjectFactory struct {
	ids      []string
	ids_lock *sync.Mutex
}

func NewARPCCentralMediaCtl() (*ARPCMediaObjectFactory, error) {
	self := new(ARPCMediaObjectFactory)
	return self, nil
}

func (self *ARPCMediaObjectFactory) genAndSaveUniqueId(
// lrc *golockerreentrancycontext.LockerReentrancyContext,
) (string, error) {
	self.ids_lock.Lock()
	defer self.ids_lock.Unlock()

	var ret string

main_loop:
	for true {
		u := uuid.NewV4()

		ret = strings.ToLower(u.String())
		for _, x := range self.ids {
			if x == ret {
				continue main_loop
			}
		}
		break
	}

	self.ids = append(self.ids, ret)

	return ret, nil
}

func (self *ARPCMediaObjectFactory) NewMediaBuffer(
	readwriteseekcloser ReadWriteSeekCloser,
	bufferInfo func(*ARPCMediaBuffer) any,
) (
	*ARPCMediaBuffer,
	error,
) {
	if readwriteseekcloser == nil {
		return nil, errors.New("readwriteseekcloser must be defined")
	}
	id, err := self.genAndSaveUniqueId()
	if err != nil {
		return nil, err
	}
	ret := new(ARPCMediaBuffer)
	ret.id = id
	ret.bufferInfo = bufferInfo
	ret.readwriteseekcloser = readwriteseekcloser

	runtime.SetFinalizer(ret, self.destroyARPCMediaBuffer)

	return ret, nil
}

func (self *ARPCMediaObjectFactory) NewMediaBufferFromBytes(
	data []byte,
	bufferInfo func(*ARPCMediaBuffer) any,
	pos int64,
	GrowOnWriteOverflow bool,
) (
	*ARPCMediaBuffer,
	error,
) {
	gif := goinmemfile.NewInMemFileFromBytes(data, pos, GrowOnWriteOverflow)
	ret, err := self.NewMediaBuffer(gif, bufferInfo)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *ARPCMediaObjectFactory) NewMediaBroadcast(
	preexisting_buffer *ARPCMediaBuffer,
	create_new_buffer bool,
	bufferInfo func(*ARPCMediaBuffer) any,
) (
	*ARPCMediaBroadcast,
	error,
) {

	var resulting_buffer *ARPCMediaBuffer

	if create_new_buffer {
		x, err := self.NewMediaBufferFromBytes(
			[]byte{},
			bufferInfo,
			0,
			true,
		)
		if err != nil {
			return nil, err
		}
		resulting_buffer = x
	} else {
		resulting_buffer = preexisting_buffer
	}

	id, err := self.genAndSaveUniqueId()
	if err != nil {
		return nil, err
	}

	ret := new(ARPCMediaBroadcast)
	ret.Buffer = resulting_buffer
	ret.id = id

	runtime.SetFinalizer(ret, self.destroyARPCMediaBroadcast)

	return ret, nil
}

func (self *ARPCMediaObjectFactory) destroyARPCMediaBuffer(o *ARPCMediaBuffer) {
	self.destroyById(o.id)
}

func (self *ARPCMediaObjectFactory) destroyARPCMediaBroadcast(o *ARPCMediaBroadcast) {
	self.destroyById(o.id)
}

func (self *ARPCMediaObjectFactory) destroyARPCMediaConn(o *ARPCMediaConn) {
	// TODO: call .Close on Conn?
	self.destroyById(o.id)
}

func (self *ARPCMediaObjectFactory) destroyById(id string) {
	self.ids_lock.Lock()
	defer self.ids_lock.Unlock()

	for i := len(self.ids) - 1; i != -1; i-- {
		if self.ids[i] == id {
			self.ids = append(self.ids[:i], self.ids[i+1:]...)
		}
	}
}
