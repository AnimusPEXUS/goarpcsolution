package goarpcsolution

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AnimusPEXUS/golockerreentrancycontext"
	"github.com/AnimusPEXUS/gouuidtools"

	"github.com/AnimusPEXUS/utils/worker"
)

// todo: find better place for this.
// but, frankly, it's the best value. which, probably doesn't need
// to be variable
const TTL_CONST_10MIN = time.Duration(time.Minute * 10)

type ARPCNodeCtlBasic struct {
	OnCallCB            func(call *ARPCCall) (error, error)
	OnUnhandledResultCB func(result *ARPCCall)

	call_id_r             *gouuidtools.UUIDRegistry
	buffer_id_r           *gouuidtools.UUIDRegistry
	transmission_id_r     *gouuidtools.UUIDRegistry
	listening_socket_id_r *gouuidtools.UUIDRegistry
	connected_socket_id_r *gouuidtools.UUIDRegistry

	calls_mtx             *sync.Mutex
	buffers_mtx           *sync.Mutex
	transmissions_mtx     *sync.Mutex
	listening_sockets_mtx *sync.Mutex
	connected_sockets_mtx *sync.Mutex

	calls             []*ARPCNodeCtlBasicCallR
	buffers           []*ARPCNodeCtlBasicBufferR
	transmissions     []*ARPCNodeCtlBasicTransmissionR
	listening_sockets []*ARPCNodeCtlBasicListeningSocketR
	connected_sockets []*ARPCNodeCtlBasicConnectedSocketR

	handlers_mtx *sync.Mutex
	handlers     []*xARPCNodeCtlBasicCallResHandlerWrapper

	wrkr01 *worker.Worker

	stop_flag bool

	node *ARPCNode

	debugName string
}

func NewARPCNodeCtlBasic() *ARPCNodeCtlBasic {
	self := new(ARPCNodeCtlBasic)
	self.debugName = "ARPCNodeCtlBasic"

	self.calls_mtx = &sync.Mutex{}
	self.buffers_mtx = &sync.Mutex{}
	self.transmissions_mtx = &sync.Mutex{}
	self.listening_sockets_mtx = &sync.Mutex{}
	self.connected_sockets_mtx = &sync.Mutex{}

	self.wrkr01 = worker.New(self.worker01)
	return self
}

func (self *ARPCNodeCtlBasic) SetNode(node *ARPCNode) {
	self.node = node
}

func (self *ARPCNodeCtlBasic) SetDebugName(name string) {
	self.debugName = fmt.Sprintf("[%s]", name)
}

func (self *ARPCNodeCtlBasic) GetDebugName() string {
	return self.debugName
}

func (self *ARPCNodeCtlBasic) DebugPrintln(data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), data...)...)
}

func (self *ARPCNodeCtlBasic) DebugPrintfln(format string, data ...any) {
	fmt.Println(append(append([]any{}, self.debugName), fmt.Sprintf(format, data...))...)
}

const (
	worker01_interval_cleanup = time.Minute
)

func (self *ARPCNodeCtlBasic) worker01(
	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),

	is_stop_flag func() bool,
) {

	lrc := new(golockerreentrancycontext.LockerReentrancyContext)

	set_starting()
	defer func() {
		for _, i := range self.calls[:] {
			self.deleteCallR_lrc(i, lrc)
		}

		for _, i := range self.buffers[:] {
			self.deleteBufferR_lrc(i, lrc)
		}

		for _, i := range self.transmissions[:] {
			self.deleteTransmissionR_lrc(i, lrc)
		}

		for _, i := range self.listening_sockets[:] {
			self.deleteListeningSocketR_lrc(i, lrc)
		}

		for _, i := range self.connected_sockets[:] {
			self.deleteConnectedSocketR_lrc(i, lrc)
		}

		set_stopped()
	}()

	var timeout_cleanup = worker01_interval_cleanup

	for true {
		if self.stop_flag {
			break
		}

		if timeout_cleanup <= 0 {
			timeout_cleanup = worker01_interval_cleanup

			for _, x := range self.calls[:] {
				if x.TTL <= 0 {
					self.deleteCallR_lrc(x, lrc)
				} else {
					x.TTL -= worker01_interval_cleanup
				}
			}

			for _, x := range self.buffers[:] {
				if x.TTL <= 0 {
					self.deleteBufferR_lrc(x, lrc)
				} else {
					x.TTL -= worker01_interval_cleanup
				}
			}

			for _, x := range self.transmissions[:] {
				if x.TTL <= 0 {
					self.deleteTransmissionR_lrc(x, lrc)
				} else {
					x.TTL -= worker01_interval_cleanup
				}
			}

			for _, x := range self.listening_sockets[:] {
				if x.TTL <= 0 {
					self.deleteListeningSocketR_lrc(x, lrc)
				} else {
					x.TTL -= worker01_interval_cleanup
				}
			}

			for _, x := range self.connected_sockets[:] {
				if x.TTL <= 0 {
					self.deleteConnectedSocketR_lrc(x, lrc)
				} else {
					x.TTL -= worker01_interval_cleanup
				}
			}

		} else {
			timeout_cleanup -= time.Second
		}

		time.Sleep(time.Second)
	}
}

func (self *ARPCNodeCtlBasic) deleteCallR_lrc(
	obj *ARPCNodeCtlBasicCallR,
	lrc *golockerreentrancycontext.LockerReentrancyContext,
) {
	lrc.LockMutex(self.calls_mtx)
	defer lrc.UnlockMutex(self.calls_mtx)

	for i := len(self.calls) - 1; i != -1; i-- {
		if self.calls[i] == obj {
			self.calls = append(self.calls[:i], self.calls[i+1:]...)
		}
	}
}

func (self *ARPCNodeCtlBasic) deleteBufferR_lrc(
	obj *ARPCNodeCtlBasicBufferR,
	lrc *golockerreentrancycontext.LockerReentrancyContext,
) {
	lrc.LockMutex(self.buffers_mtx)
	defer lrc.UnlockMutex(self.buffers_mtx)

	for i := len(self.buffers) - 1; i != -1; i-- {
		if self.buffers[i] == obj {
			self.buffers = append(self.buffers[:i], self.buffers[i+1:]...)
		}
	}
}

func (self *ARPCNodeCtlBasic) deleteTransmissionR_lrc(
	obj *ARPCNodeCtlBasicTransmissionR,
	lrc *golockerreentrancycontext.LockerReentrancyContext,
) {
	lrc.LockMutex(self.transmissions_mtx)
	defer lrc.UnlockMutex(self.transmissions_mtx)

	for i := len(self.transmissions) - 1; i != -1; i-- {
		if self.transmissions[i] == obj {
			self.transmissions = append(self.transmissions[:i], self.transmissions[i+1:]...)
		}
	}
}

func (self *ARPCNodeCtlBasic) deleteListeningSocketR_lrc(
	obj *ARPCNodeCtlBasicListeningSocketR,
	lrc *golockerreentrancycontext.LockerReentrancyContext,
) {
	lrc.LockMutex(self.listening_sockets_mtx)
	defer lrc.UnlockMutex(self.listening_sockets_mtx)

	for i := len(self.listening_sockets) - 1; i != -1; i-- {
		if self.listening_sockets[i] == obj {
			self.listening_sockets = append(
				self.listening_sockets[:i],
				self.listening_sockets[i+1:]...,
			)
		}
	}
}

func (self *ARPCNodeCtlBasic) deleteConnectedSocketR_lrc(
	obj *ARPCNodeCtlBasicConnectedSocketR,
	lrc *golockerreentrancycontext.LockerReentrancyContext,
) {
	lrc.LockMutex(self.connected_sockets_mtx)
	defer lrc.UnlockMutex(self.connected_sockets_mtx)

	for i := len(self.connected_sockets) - 1; i != -1; i-- {
		if self.connected_sockets[i] == obj {
			self.connected_sockets = append(
				self.connected_sockets[:i],
				self.connected_sockets[i+1:]...,
			)
		}
	}
}

// Call and Reply essentially the same,
// but Call has reply_to_id field set to nil, and Reply doesn't use
// 'name' field

func (self *ARPCNodeCtlBasic) Call(
	name string,
	args []*ARPCCallArg,

	unhandled bool,
	rh *ARPCNodeCtlBasicCallResHandler,
	response_timeout time.Duration,
	// note: not sure if this is needed here. but it's possible
	//       should be used to retreive call_id before actual
	//       call is done
	// request_id_hook *ARPCNodeCtlBasicNewCallIdHook,
) (ret_any *gouuidtools.UUID, ret_err error) {

	call_id, err := self.call_id_r.GenUUID()
	if err != nil {
		return
	}

	err = self.saveCall(
		call_id,
		nil,
		name,
		args,
		TTL_CONST_10MIN,
	)

	if err != nil {
		return nil, err
	}

	err = self.node.NewCall(
		call_id,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return call_id, nil
}

func (self *ARPCNodeCtlBasic) Reply(
	reply_to_id *gouuidtools.UUID,
	args ...*ARPCCallArg,
) (
	err error,
) {
	call_id, err := self.call_id_r.GenUUID()
	if err != nil {
		return
	}

	err = self.saveCall(
		call_id,
		reply_to_id,
		"",
		args,
		TTL_CONST_10MIN,
	)

	if err != nil {
		return err
	}

	err = self.node.NewCall(
		call_id,
		reply_to_id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (self *ARPCNodeCtlBasic) saveCall(
	call_id *gouuidtools.UUID,
	reply_to_id *gouuidtools.UUID,

	name string,
	args []*ARPCCallArg,

	ttl time.Duration,
) error {

	lrc := new(golockerreentrancycontext.LockerReentrancyContext)

	if (name != "" && reply_to_id != nil) ||
		(name == "" && reply_to_id == nil) {
		return errors.New(
			"only 'name' or only 'reply_to_id' must be set",
		)
	}

	var err error

	for _, i := range args {
		if i == nil {
			return errors.New("nil is not acceptable amont args")
		}

		if err = i.IsValidError(); err != nil {
			return err
		}
	}

	buffer_w := make([]*ARPCNodeCtlBasicBufferR, 0)
	transmission_w := make([]*ARPCNodeCtlBasicTransmissionR, 0)
	listening_socket_w := make([]*ARPCNodeCtlBasicListeningSocketR, 0)
	connected_socket_w := make([]*ARPCNodeCtlBasicConnectedSocketR, 0)

	for _, i := range args {
		if i.Buffer != nil {
			uuid := i.Buffer.Id
			if uuid == nil || uuid.IsNil() {
				uuid, err = self.buffer_id_r.GenUUID()
				if err != nil {
					return err
				}
			}

			b := &ARPCNodeCtlBasicBufferR{
				Ctl:      self,
				BufferId: uuid,
				Buffer:   i.Buffer.Payload,
				TTL:      TTL_CONST_10MIN,
			}

			buffer_w = append(buffer_w, b)
		}
	}

	for _, i := range args {
		if i.Transmission != nil {
			uuid := i.Transmission.Id
			if uuid == nil || uuid.IsNil() {
				uuid, err = self.transmission_id_r.GenUUID()
				if err != nil {
					return err
				}
			}

			b := &ARPCNodeCtlBasicTransmissionR{
				Ctl:            self,
				TransmissionId: uuid,
				Transmission:   i.Transmission.Payload,
				TTL:            TTL_CONST_10MIN,
			}

			transmission_w = append(transmission_w, b)
		}
	}

	for _, i := range args {
		if i.ListeningSocket != nil {
			uuid := i.ListeningSocket.Id
			if uuid == nil || uuid.IsNil() {
				uuid, err = self.listening_socket_id_r.GenUUID()
				if err != nil {
					return err
				}
			}

			b := &ARPCNodeCtlBasicListeningSocketR{
				Ctl:               self,
				ListeningSocketId: uuid,
				ListeningSocket:   i.ListeningSocket.Payload,
				TTL:               TTL_CONST_10MIN,
			}

			listening_socket_w = append(listening_socket_w, b)
		}
	}

	for _, i := range args {
		if i.ConnectedSocket != nil {
			uuid := i.ConnectedSocket.Id
			if uuid == nil || uuid.IsNil() {
				uuid, err = self.connected_socket_id_r.GenUUID()
				if err != nil {
					return err
				}
			}

			b := &ARPCNodeCtlBasicConnectedSocketR{
				Ctl:               self,
				ConnectedSocketId: uuid,
				ConnectedSocket:   i.ConnectedSocket.Payload,
				TTL:               TTL_CONST_10MIN,
			}

			connected_socket_w = append(connected_socket_w, b)
		}
	}

	call := &ARPCNodeCtlBasicCallR{
		Ctl:       self,
		CallId:    call_id,
		ReplyToId: reply_to_id,
		Name:      name,
		Args:      args,
		TTL:       TTL_CONST_10MIN,
	}

	lrc.LockMutex(self.calls_mtx)
	defer lrc.UnlockMutex(self.calls_mtx)

	lrc.LockMutex(self.buffers_mtx)
	defer lrc.UnlockMutex(self.buffers_mtx)

	lrc.LockMutex(self.transmissions_mtx)
	defer lrc.UnlockMutex(self.transmissions_mtx)

	lrc.LockMutex(self.listening_sockets_mtx)
	defer lrc.UnlockMutex(self.listening_sockets_mtx)

	lrc.LockMutex(self.connected_sockets_mtx)
	defer lrc.UnlockMutex(self.connected_sockets_mtx)

	self.calls = append(self.calls, call)

	self.buffers = append(
		self.buffers,
		buffer_w...,
	)
	self.transmissions = append(
		self.transmissions,
		transmission_w...,
	)
	self.listening_sockets = append(
		self.listening_sockets,
		listening_socket_w...,
	)
	self.connected_sockets = append(
		self.connected_sockets,
		connected_socket_w...,
	)

	return nil
}

// 'R' at the end of next structs - stands for 'Record'

type ARPCNodeCtlBasicCallR struct {
	Ctl *ARPCNodeCtlBasic

	CallId    *gouuidtools.UUID
	ReplyToId *gouuidtools.UUID

	Name string
	Args []*ARPCCallArg

	TTL time.Duration
}

func (self *ARPCNodeCtlBasicCallR) Deleted() {

}

type ARPCNodeCtlBasicBufferR struct {
	Ctl *ARPCNodeCtlBasic

	BufferId *gouuidtools.UUID

	Buffer ARPCBufferI

	TTL time.Duration
}

func (self *ARPCNodeCtlBasicBufferR) Deleted() {

}

type ARPCNodeCtlBasicTransmissionR struct {
	Ctl *ARPCNodeCtlBasic

	TransmissionId *gouuidtools.UUID

	Transmission ARPCTransmissionI

	TTL time.Duration
}

func (self *ARPCNodeCtlBasicTransmissionR) Deleted() {

}

type ARPCNodeCtlBasicListeningSocketR struct {
	Ctl *ARPCNodeCtlBasic

	ListeningSocketId *gouuidtools.UUID

	ListeningSocket ARPCListeningSocketI

	TTL time.Duration
}

func (self *ARPCNodeCtlBasicListeningSocketR) Deleted() {

}

type ARPCNodeCtlBasicConnectedSocketR struct {
	Ctl *ARPCNodeCtlBasic

	ConnectedSocketId *gouuidtools.UUID

	ConnectedSocket ARPCConnectedSocketI

	TTL time.Duration
}

func (self *ARPCNodeCtlBasicConnectedSocketR) Deleted() {

}

type xARPCNodeCtlBasicCallResHandlerWrapper struct {
	handler *ARPCNodeCtlBasicCallResHandler
	id      *gouuidtools.UUID
	timeout time.Duration
}

type ARPCNodeCtlBasicCallResHandler struct {
	OnTimeout  func()
	OnClose    func()
	OnResponse func(args *ARPCCall)
}
